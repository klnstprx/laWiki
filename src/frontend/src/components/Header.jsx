import { useRef, useState } from "react";
import {
  AppBar,
  Box,
  Button,
  CircularProgress,
  ClickAwayListener,
  InputBase,
  Paper,
  Popper,
  Toolbar,
  Typography,
  useMediaQuery,
  useTheme,
} from "@mui/material";
import { Home as HomeIcon, Search as SearchIcon } from "@mui/icons-material";
import { alpha, styled } from "@mui/material/styles";
import { useNavigate } from "react-router-dom";
import { searchWikis } from "../api/WikiApi";
import { searchEntries } from "../api/EntryApi";
import { searchComments } from "../api/CommentApi";
import { searchVersions } from "../api/VersionApi";
import { getUsersByIds } from "../api/AuthApi";
import SearchResultsList from "./SearchResultsList";
import LoginButton from "./LoginButton";

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
  [theme.breakpoints.up("md")]: {
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
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down("md"));
  const [showSearch, setShowSearch] = useState(false);

  const handleSearchToggle = () => {
    setShowSearch(!showSearch);
  };

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
      // Perform searches
      const [wikis, entries, comments, versions] = await Promise.all([
        searchWikis({ title: query }),
        searchEntries({ title: query }),
        searchComments({ content: query }),
        searchVersions({ content: query }),
      ]);

      // Ensure that each result is an array
      const wikisArray = Array.isArray(wikis) ? wikis : [];
      const entriesArray = Array.isArray(entries) ? entries : [];
      const commentsArray = Array.isArray(comments) ? comments : [];
      const versionsArray = Array.isArray(versions) ? versions : [];

      // Map to hold user IDs and names for display purposes
      const userIdToName = {};

      // Collect user IDs from search results
      const allUserIds = new Set();
      entriesArray.forEach((entry) => {
        if (entry.author) allUserIds.add(entry.author);
      });
      commentsArray.forEach((comment) => {
        if (comment.author) allUserIds.add(comment.author);
      });
      versionsArray.forEach((version) => {
        if (version.editor) allUserIds.add(version.editor);
      });

      // Fetch user data for IDs
      const missingUserIds = Array.from(allUserIds).filter(
        (id) => !userIdToName[id],
      );
      if (missingUserIds.length > 0) {
        const additionalUsers = await getUsersByIds(missingUserIds);
        additionalUsers.forEach((user) => {
          userIdToName[user.id] = user.name;
        });
      }

      // Prepare the results with user names
      const processedEntries = entriesArray.map((entry) => ({
        ...entry,
        authorName: userIdToName[entry.author] || "Usuario desconocido",
      }));

      const processedComments = commentsArray.map((comment) => ({
        ...comment,
        authorName: userIdToName[comment.author] || "Usuario desconocido",
      }));

      const processedVersions = versionsArray.map((version) => ({
        ...version,
        editorName: userIdToName[version.editor] || "Usuario desconocido",
      }));

      // Set the search results
      setSearchResults({
        wikis: wikisArray,
        entries: processedEntries,
        comments: processedComments,
        versions: processedVersions,
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

  const goHome = () => {
    navigate("/");
  };
  return (
    <AppBar position="static">
      <Toolbar sx={{ gap: 2 }}>
        {isMobile
          ? (
            <Button
              variant="contained"
              color="secondary"
              onClick={goHome}
              sx={{
                p: 1,
                minWidth: "auto",
              }}
            >
              <HomeIcon />
            </Button>
          )
          : (
            <Button
              variant="contained"
              color="secondary"
              sx={{
                p: 1,
                minWidth: "auto",
              }}
              onClick={goHome}
              startIcon={<HomeIcon />}
            >
              <Typography
                variant="h6"
                noWrap
                sx={{
                  color: "inherit",
                  textDecoration: "none",
                }}
              >
                Home
              </Typography>
            </Button>
          )}
        <Box sx={{ flexGrow: 1 }} />
        {isMobile
          ? (
            <>
              <Button
                variant="contained"
                color="secondary"
                sx={{
                  p: 1,
                  minWidth: "auto",
                }}
                onClick={handleSearchToggle}
              >
                <SearchIcon />
              </Button>
            </>
          )
          : (
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
          )}

        <LoginButton />
      </Toolbar>
      {isMobile && showSearch && (
        <Box sx={{ p: 3 }}>
          <Search sx={{ m: 0 }} ref={searchRef}>
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
                        <Typography variant="body2">
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
        </Box>
      )}
    </AppBar>
  );
};

export default Header;
