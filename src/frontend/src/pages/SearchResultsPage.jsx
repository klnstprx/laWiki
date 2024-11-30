import { useEffect, useState } from "react";
import { useLocation } from "react-router-dom";
import {
  Container,
  Typography,
  CircularProgress,
  Grid2,
  Card,
  CardContent,
} from "@mui/material";

// Import your API functions
import { searchWikis } from "../api/WikiApi";
import { searchEntries } from "../api/EntryApi";
import { searchComments } from "../api/CommentApi";
import { searchVersions } from "../api/VersionApi";

const SearchResultsPage = () => {
  const location = useLocation();
  const queryParams = new URLSearchParams(location.search);
  const [results, setResults] = useState({
    wikis: [],
    entries: [],
    comments: [],
    versions: [],
  });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    const performSearch = async () => {
      setLoading(true);
      setError("");

      const queries = {};

      // Extract parameters for each service

      // For basic search (from header)
      const q = queryParams.get("q");
      if (q) {
        queries.wikiParams = { title: q };
        queries.entryParams = { title: q };
        queries.commentParams = { content: q };
        queries.versionParams = { content: q };
      }

      // For advanced search
      if (queryParams.get("wikiTitle") || queryParams.get("wikiCategory")) {
        queries.wikiParams = {
          ...(queries.wikiParams || {}),
          ...(queryParams.get("wikiTitle")
            ? { title: queryParams.get("wikiTitle") }
            : {}),
          ...(queryParams.get("wikiCategory")
            ? { category: queryParams.get("wikiCategory") }
            : {}),
        };
      }

      if (queryParams.get("entryTitle") || queryParams.get("entryAuthor")) {
        queries.entryParams = {
          ...(queries.entryParams || {}),
          ...(queryParams.get("entryTitle")
            ? { title: queryParams.get("entryTitle") }
            : {}),
          ...(queryParams.get("entryAuthor")
            ? { author: queryParams.get("entryAuthor") }
            : {}),
        };
      }

      if (
        queryParams.get("commentContent") ||
        queryParams.get("commentAuthor")
      ) {
        queries.commentParams = {
          ...(queries.commentParams || {}),
          ...(queryParams.get("commentContent")
            ? { content: queryParams.get("commentContent") }
            : {}),
          ...(queryParams.get("commentAuthor")
            ? { author: queryParams.get("commentAuthor") }
            : {}),
        };
      }

      if (
        queryParams.get("versionContent") ||
        queryParams.get("versionEditor")
      ) {
        queries.versionParams = {
          ...(queries.versionParams || {}),
          ...(queryParams.get("versionContent")
            ? { content: queryParams.get("versionContent") }
            : {}),
          ...(queryParams.get("versionEditor")
            ? { editor: queryParams.get("versionEditor") }
            : {}),
        };
      }

      try {
        const [wikis, entries, comments, versions] = await Promise.all([
          queries.wikiParams
            ? searchWikis(queries.wikiParams)
            : Promise.resolve([]),
          queries.entryParams
            ? searchEntries(queries.entryParams)
            : Promise.resolve([]),
          queries.commentParams
            ? searchComments(queries.commentParams)
            : Promise.resolve([]),
          queries.versionParams
            ? searchVersions(queries.versionParams)
            : Promise.resolve([]),
        ]);

        setResults({
          wikis: wikis || [],
          entries: entries || [],
          comments: comments || [],
          versions: versions || [],
        });
      } catch (err) {
        setError(err.message || "Error during search");
      } finally {
        setLoading(false);
      }
    };

    performSearch();
  }, [location.search]);

  if (loading) {
    return (
      <Container sx={{ mt: 4, textAlign: "center" }}>
        <CircularProgress />
        <Typography variant="body1">Loading results...</Typography>
      </Container>
    );
  }

  if (error) {
    return (
      <Container sx={{ mt: 4 }}>
        <Typography variant="h6" color="error">
          Error: {error}
        </Typography>
      </Container>
    );
  }

  return (
    <Container maxWidth="md" sx={{ mt: 4 }}>
      <Typography variant="h4" gutterBottom>
        Search Results
      </Typography>

      {/* No Results */}
      {results.wikis.length === 0 &&
      results.entries.length === 0 &&
      results.comments.length === 0 &&
      results.versions.length === 0 ? (
        <Typography variant="h6">No results found.</Typography>
      ) : null}

      {/* Wikis */}
      {results.wikis.length > 0 && (
        <div>
          <Typography variant="h5" gutterBottom>
            Wikis
          </Typography>
          <Grid2 container spacing={2}>
            {results.wikis.map((wiki) => (
              <Grid2 xs={12} key={wiki.id}>
                <Card>
                  <CardContent>
                    <Typography variant="h6">{wiki.title}</Typography>
                    <Typography variant="body2">{wiki.description}</Typography>
                    <Typography variant="caption">
                      Category: {wiki.category}
                    </Typography>
                  </CardContent>
                </Card>
              </Grid2>
            ))}
          </Grid2>
        </div>
      )}

      {/* Entries */}
      {results.entries.length > 0 && (
        <div>
          <Typography variant="h5" gutterBottom sx={{ mt: 4 }}>
            Entries
          </Typography>
          <Grid2 container spacing={2}>
            {results.entries.map((entry) => (
              <Grid2 xs={12} key={entry.id}>
                <Card>
                  <CardContent>
                    <Typography variant="h6">{entry.title}</Typography>
                    <Typography variant="body2">
                      Author: {entry.author}
                    </Typography>
                    <Typography variant="caption">
                      Wiki ID: {entry.wiki_id}
                    </Typography>
                  </CardContent>
                </Card>
              </Grid2>
            ))}
          </Grid2>
        </div>
      )}

      {/* Comments */}
      {results.comments.length > 0 && (
        <div>
          <Typography variant="h5" gutterBottom sx={{ mt: 4 }}>
            Comments
          </Typography>
          <Grid2 container spacing={2}>
            {results.comments.map((comment) => (
              <Grid2 xs={12} key={comment.id}>
                <Card>
                  <CardContent>
                    <Typography variant="h6">
                      Author: {comment.author}
                    </Typography>
                    <Typography variant="body2">{comment.content}</Typography>
                    <Typography variant="caption">
                      Rating: {comment.rating}
                    </Typography>
                  </CardContent>
                </Card>
              </Grid2>
            ))}
          </Grid2>
        </div>
      )}

      {/* Versions */}
      {results.versions.length > 0 && (
        <div>
          <Typography variant="h5" gutterBottom sx={{ mt: 4 }}>
            Versions
          </Typography>
          <Grid2 container spacing={2}>
            {results.versions.map((version) => (
              <Grid2 xs={12} key={version.id}>
                <Card>
                  <CardContent>
                    <Typography variant="h6">
                      Editor: {version.editor}
                    </Typography>
                    <Typography variant="body2">{version.content}</Typography>
                  </CardContent>
                </Card>
              </Grid2>
            ))}
          </Grid2>
        </div>
      )}
    </Container>
  );
};

export default SearchResultsPage;
