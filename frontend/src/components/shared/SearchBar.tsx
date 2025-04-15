import React, { useState, useEffect, useCallback } from 'react';
import {
  Box,
  TextField,
  Autocomplete,
  IconButton,
  InputAdornment,
  CircularProgress,
  Typography,
  AutocompleteRenderInputParams,
  TextFieldProps,
} from '@mui/material';
import {
  Search as SearchIcon,
  Clear as ClearIcon,
} from '@mui/icons-material';
import debounce from 'lodash/debounce';
import { SearchSuggestion } from '../../types/search';
import { searchApi } from '../../services/search';

interface SearchBarProps {
  onSearch: (query: string) => void;
  placeholder?: string;
  fullWidth?: boolean;
}

export const SearchBar: React.FC<SearchBarProps> = ({
  onSearch,
  placeholder = 'Search services...',
  fullWidth = false,
}) => {
  const [query, setQuery] = useState('');
  const [suggestions, setSuggestions] = useState<SearchSuggestion[]>([]);
  const [loading, setLoading] = useState(false);
  const [open, setOpen] = useState(false);

  // Debounced function to fetch suggestions
  const debouncedFetchSuggestions = useCallback(
    debounce(async (searchQuery: string) => {
      if (!searchQuery) {
        setSuggestions([]);
        return;
      }

      try {
        setLoading(true);
        const response = await searchApi.getSearchSuggestions(searchQuery);
        if (response.success) {
          setSuggestions(response.data);
        } else {
          console.error('Failed to fetch suggestions:', response.message);
          setSuggestions([]);
        }
      } catch (error) {
        console.error('Failed to fetch suggestions:', error);
        setSuggestions([]);
      } finally {
        setLoading(false);
      }
    }, 300),
    []
  );

  useEffect(() => {
    debouncedFetchSuggestions(query);
  }, [query, debouncedFetchSuggestions]);

  const handleSearch = () => {
    onSearch(query);
  };

  const handleKeyPress = (event: React.KeyboardEvent) => {
    if (event.key === 'Enter') {
      handleSearch();
    }
  };

  const handleClear = () => {
    setQuery('');
    setSuggestions([]);
    onSearch('');
  };

  const renderInput = (params: AutocompleteRenderInputParams) => {
    const { InputProps, InputLabelProps, ...rest } = params;

    return (
      <TextField
        {...rest}
        placeholder={placeholder}
        variant="outlined"
        size="medium"
        InputLabelProps={{
          ...InputLabelProps,
          className: InputLabelProps?.className || '',
          style: InputLabelProps?.style || {},
        }}
        InputProps={{
          ...InputProps,
          startAdornment: (
            <InputAdornment position="start">
              <SearchIcon color="action" />
            </InputAdornment>
          ),
          endAdornment: (
            <React.Fragment>
              {loading ? (
                <CircularProgress color="inherit" size={20} />
              ) : query ? (
                <IconButton
                  size="small"
                  onClick={handleClear}
                  sx={{ mr: 1 }}
                >
                  <ClearIcon />
                </IconButton>
              ) : null}
              {InputProps.endAdornment}
            </React.Fragment>
          ),
        }}
        onKeyPress={handleKeyPress}
      />
    );
  };

  return (
    <Box sx={{ width: fullWidth ? '100%' : 'auto' }}>
      <Autocomplete
        open={open}
        onOpen={() => setOpen(true)}
        onClose={() => setOpen(false)}
        freeSolo
        options={suggestions}
        getOptionLabel={(option): string => 
          typeof option === 'string' ? option : option.label
        }
        filterOptions={(x) => x}
        inputValue={query}
        onInputChange={(_, newValue) => {
          setQuery(newValue);
        }}
        onChange={(_, newValue) => {
          const value = typeof newValue === 'string' ? newValue : newValue?.value || '';
          setQuery(value);
          onSearch(value);
        }}
        renderInput={renderInput}
        renderOption={(props, option) => (
          <Box
            component="li"
            {...props}
            key={option.value}
            sx={{
              p: 1,
              cursor: 'pointer',
              '&:hover': {
                backgroundColor: 'action.hover',
              },
            }}
          >
            <Typography>{option.label}</Typography>
          </Box>
        )}
      />
    </Box>
  );
};

export default SearchBar;
