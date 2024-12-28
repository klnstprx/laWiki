import { useRef, useState } from "react";
import {
  AppBar,
  Button,
  CircularProgress,
  ClickAwayListener,
  IconButton,
  InputBase,
  Paper,
  Popper,
  Toolbar,
  Typography,
} from "@mui/material";
import { Home as HomeIcon, Search as SearchIcon } from "@mui/icons-material";
import { alpha, styled } from "@mui/material/styles";
import { Link as RouterLink, useNavigate } from "react-router-dom";
import { searchWikis } from "../api/WikiApi";
import { searchEntries } from "../api/EntryApi";
import { searchComments } from "../api/CommentApi";
import { searchVersions } from "../api/VersionApi";
import SearchResultsList from "./SearchResultsList";
import UserNav from "./UserNav";

const Search = styled("div")(({ theme }) => ({
  position: "relative",
  borderRadius: theme.shape.borderRadius,
  backgroundColor: alpha(theme.palette.common.white, 0.15),
  "&:hover": {
    backgroundColor: alpha(theme.palette.common.white, 0.25),
  },
  marginLeft: theme.spacing(2),
  marginRight: theme.spacing(2),
  width: "100%",
  [theme.breakpoints.up("sm")]: {
    marginLeft: theme.spacing(2),
    width: "auto",
  },
}));

const SearchIconWrapper = styled("div")(({ theme }) => ({
  padding: theme.spacing(0, 1),
  height: "100%",
  position: "absolute",
  pointerEvents: "none",
  display: "flex",
  alignItems: "center",
  justifyContent: "center",
}));

const StyledInputBase = styled(InputBase)(({ theme }) => ({
  color: "inherit",
  width: "100%",
  "& .MuiInputBase-input": {
    padding: theme.spacing(1, 1, 1, 0),
    // vertical padding + font size from searchIcon
    paddingLeft: `calc(1em + ${theme.spacing(4)})`,
    transition: theme.transitions.create("width"),
    width: "12ch",
    [theme.breakpoints.up("md")]: {
      width: "20ch",
    },
  },
}));

const Header = () => {
  const navigate = useNavigate();
  const [searchQuery, setSearchQuery] = useState("");

  const [searchResults, setSearchResults] = useState({
    wikis: [],
    entries: [],
  });
  const [open, setOpen] = useState(false);
  const [loading, setLoading] = useState(false);

  const searchRef = useRef(null);
  const debounceTimeoutRef = useRef(null);

  const handleSearchInputChange = (event) => {
    const query = event.target.value;
    setSearchQuery(query);

    // Clear previous debounce timeout
    if (debounceTimeoutRef.current) {
      clearTimeout(debounceTimeoutRef.current);
    }

    // Set new debounce timeout
    debounceTimeoutRef.current = setTimeout(() => {
      if (query.trim()) {
        performSearch(query.trim());
      } else {
        setSearchResults({ wikis: [], entries: [] });
        setOpen(false);
      }
    }, 300); // 300ms debounce
  };

  const performSearch = async (query) => {
    setLoading(true);
    try {
      const [wikis, entries, comments, versions] = await Promise.all([
        searchWikis({ title: query }),
        searchEntries({ title: query }),
        searchComments({ content: query }),
        searchVersions({ content: query }),
      ]);

      setSearchResults({
        wikis: Array.isArray(wikis) ? wikis : [],
        entries: Array.isArray(entries) ? entries : [],
        comments: Array.isArray(comments) ? comments : [],
        versions: Array.isArray(versions) ? versions : [],
      });
      setOpen(true);
    } catch (error) {
      console.error("Error during search:", error);
      setSearchResults({
        wikis: [],
        entries: [],
        comments: [],
        versions: [],
      });
      setOpen(false);
    } finally {
      setLoading(false);
    }
  };

  const handleClickAway = () => {
    setOpen(false);
  };

  const handleAdvancedSearch = () => {
    navigate("/advanced-search");
  };

  return (
    <AppBar position="static">
      <Toolbar>
        <IconButton
          component={RouterLink}
          to="/"
          sx={{
            display: { xs: "block", sm: "none" },
            color: "inherit",
            textDecoration: "none",
            mr: 2,
          }}
        >
          <HomeIcon />
        </IconButton>

        <Typography
          variant="h6"
          noWrap
          component={RouterLink}
          to="/"
          sx={{
            display: { xs: "none", sm: "block" },
            color: "inherit",
            textDecoration: "none",
            flexGrow: 1,
            "&:hover": {
              color: "pink",
            },
          }}
        >
          Home
        </Typography>

        <UserNav />

        <Search ref={searchRef}>
          <SearchIconWrapper>
            <SearchIcon />
          </SearchIconWrapper>
          <StyledInputBase
            placeholder="Buscar…"
            inputProps={{ "aria-label": "search" }}
            value={searchQuery}
            onChange={handleSearchInputChange}
          />
          <Popper
            open={open}
            anchorEl={searchRef.current}
            placement="bottom-start"
            style={{ zIndex: 1100 }}
          >
            <ClickAwayListener onClickAway={handleClickAway}>
              <Paper
                elevation={3}
                style={{
                  width: searchRef.current
                    ? searchRef.current.clientWidth
                    : 200,
                  maxHeight: "400px",
                  overflowY: "auto",
                }}
              >
                {loading
                  ? (
                    <div
                      style={{
                        display: "flex",
                        alignItems: "center",
                        padding: "10px",
                      }}
                    >
                      <CircularProgress size={24} />
                      <Typography variant="body2" sx={{ ml: 2 }}>
                        Cargando...
                      </Typography>
                    </div>
                  )
                  : (
                    <SearchResultsList
                      results={searchResults}
                      onItemClick={handleClickAway}
                    />
                  )}
              </Paper>
            </ClickAwayListener>
          </Popper>
        </Search>

        <Button
          variant="contained"
          color="secondary"
          onClick={handleAdvancedSearch}
          startIcon={<SearchIcon />}
        >
          Búsqueda Avanzada
        </Button>
      </Toolbar>
    </AppBar>
  );
};

export default Header;
