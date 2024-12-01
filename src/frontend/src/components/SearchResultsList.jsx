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

const SearchResultsList = ({
  results = { wikis: [], entries: [] },
  onItemClick,
}) => {
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

SearchResultsList.propTypes = {
  results: PropTypes.shape({
    wikis: PropTypes.array,
    entries: PropTypes.array,
  }),
  onItemClick: PropTypes.func,
};

export default SearchResultsList;
