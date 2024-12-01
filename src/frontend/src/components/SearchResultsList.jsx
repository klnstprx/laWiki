import PropTypes from "prop-types";
import {
  Typography,
  List,
  ListItem,
  ListItemText,
  Divider,
  ListItemButton,
} from "@mui/material";
import { Link as RouterLink } from "react-router-dom";

const SearchResultsList = ({ results = {}, onItemClick }) => {
  const { wikis = [], entries = [], comments = [], versions = [] } = results;

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
                  onClick={onItemClick}
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
                  onClick={onItemClick}
                >
                  <ListItemText
                    primary={entry.title}
                    secondary={`Autor: ${entry.author}`}
                  />
                </ListItemButton>
              </ListItem>
            ))}
          </List>
          <Divider />
        </>
      )}

      {comments.length > 0 && (
        <>
          <Typography variant="subtitle1" sx={{ pl: 2, pt: 1 }}>
            Comentarios
          </Typography>
          <List>
            {comments.map((comment) => (
              <ListItem key={comment.id} disablePadding>
                <ListItemButton
                  component={RouterLink}
                  to={`/version/${comment.entry_id}/${comment.version_id}`}
                  onClick={onItemClick}
                >
                  <ListItemText
                    primary={comment.content}
                    secondary={`Autor: ${comment.author}`}
                  />
                </ListItemButton>
              </ListItem>
            ))}
          </List>
          <Divider />
        </>
      )}

      {versions.length > 0 && (
        <>
          <Typography variant="subtitle1" sx={{ pl: 2, pt: 1 }}>
            Versiones
          </Typography>
          <List>
            {versions.map((version) => (
              <ListItem key={version.id} disablePadding>
                <ListItemButton
                  component={RouterLink}
                  to={`/entrada/${version.entry_id}/${version.id}`}
                  onClick={onItemClick}
                >
                  <ListItemText
                    primary={`Editor: ${version.editor}`}
                    secondary={`Fecha: ${new Date(version.created_at).toLocaleDateString()}`}
                  />
                </ListItemButton>
              </ListItem>
            ))}
          </List>
          <Divider />
        </>
      )}

      {wikis.length === 0 &&
        entries.length === 0 &&
        comments.length === 0 &&
        versions.length === 0 && (
          <Typography variant="body2" sx={{ pl: 2, pt: 1 }}>
            No se encontraron resultados.
          </Typography>
        )}
    </>
  );
};

SearchResultsList.propTypes = {
  results: PropTypes.shape({
    wikis: PropTypes.array,
    entries: PropTypes.array,
    comments: PropTypes.array,
    versions: PropTypes.array,
  }),
  onItemClick: PropTypes.func,
};

export default SearchResultsList;
