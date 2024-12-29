import { useState } from "react";
import PropTypes from "prop-types";
import { LanguageContext } from "./LanguageContext";

export const LanguageProvider = ({ children }) => {
  const [selectedOption, setSelectedOption] = useState(null);

  return (
    <LanguageContext.Provider value={{ selectedOption, setSelectedOption }}>
      {children}
    </LanguageContext.Provider>
  );
};

LanguageProvider.propTypes = {
  children: PropTypes.node.isRequired,
};