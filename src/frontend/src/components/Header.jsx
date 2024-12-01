import { useState, useRef } from "react";
import {
  AppBar,
  Toolbar,
  Typography,
  InputBase,
  IconButton,
  ListItemButton,
  Button,
  Popper,
  Paper,
  ClickAwayListener,
  List,
  ListItem,
  ListItemText,
  Divider,
  CircularProgress,
} from "@mui/material";
import { Search as SearchIcon, Home as HomeIcon } from "@mui/icons-material";
import { alpha, styled } from "@mui/material/styles";
import { Link as RouterLink, useNavigate } from "react-router-dom";
import { searchWikis } from "../api/WikiApi";
import { searchEntries } from "../api/EntryApi";

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
      const [wikis, entries] = await Promise.all([
        searchWikis({ title: query }),
        searchEntries({ title: query }),
      ]);

      // Ensure wikis and entries are arrays
      setSearchResults({
        wikis: Array.isArray(wikis) ? wikis : [],
        entries: Array.isArray(entries) ? entries : [],
      });
      setOpen(true);
    } catch (error) {
      console.error("Error during search:", error);
      setSearchResults({ wikis: [], entries: [] });
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
          }}
        >
          Home
        </Typography>

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
                {loading ? (
                  <ListItem>
                    <CircularProgress size={24} />
                    <ListItemText primary="Cargando..." />
                  </ListItem>
                ) : (
                  renderSearchResults(searchResults, handleClickAway)
                )}
              </Paper>
            </ClickAwayListener>
          </Popper>
        </Search>

        <Button color="inherit" onClick={handleAdvancedSearch}>
          Búsqueda Avanzada
        </Button>
      </Toolbar>
    </AppBar>
  );
};

const renderSearchResults = (
  results = { wikis: [], entries: [] },
  handleClose,
) => {
  const wikis = Array.isArray(results.wikis) ? results.wikis : [];
  const entries = Array.isArray(results.entries) ? results.entries : [];

  return (
    <>
      {wikis.length > 0 && (
        <>
          <Typography variant="subtitle1" sx={{ pl: 2, pt: 1 }}>
            Wikis
          </Typography>
          <List>
            {wikis.map((wiki) => (
              <ListItem key={wiki.id} disablePadding>
                <ListItemButton
                  component={RouterLink}
                  to={`/wiki/${wiki.id}`}
                  onClick={handleClose}
                >
                  <ListItemText
                    primary={wiki.title}
                    secondary={wiki.description}
                  />
                </ListItemButton>
              </ListItem>
            ))}
          </List>
          <Divider />
        </>
      )}
      {entries.length > 0 && (
        <>
          <Typography variant="subtitle1" sx={{ pl: 2, pt: 1 }}>
            Entradas
          </Typography>
          <List>
            {entries.map((entry) => (
              <ListItem key={entry.id} disablePadding>
                <ListItemButton
                  component={RouterLink}
                  to={`/entrada/${entry.id}`}
                  onClick={handleClose}
                >
                  <ListItemText
                    primary={entry.title}
                    secondary={`Autor: ${entry.author}`}
                  />
                </ListItemButton>
              </ListItem>
            ))}
          </List>
        </>
      )}
      {wikis.length === 0 && entries.length === 0 && (
        <Typography variant="body2" sx={{ pl: 2, pt: 1 }}>
          No se encontraron resultados.
        </Typography>
      )}
    </>
  );
};

export default Header;
