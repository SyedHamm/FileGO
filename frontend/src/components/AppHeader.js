import React from 'react';
import AppBar from '@mui/material/AppBar';
import Toolbar from '@mui/material/Toolbar';
import Typography from '@mui/material/Typography';
import IconButton from '@mui/material/IconButton';
import Brightness4Icon from '@mui/icons-material/Brightness4';
import Brightness7Icon from '@mui/icons-material/Brightness7';
import StorageIcon from '@mui/icons-material/Storage';

export const AppHeader = ({ darkMode, onThemeChange }) => {
  return (
    <AppBar position="static">
      <Toolbar>
        <StorageIcon sx={{ mr: 2 }} />
        <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
          Distributed File System
        </Typography>
        <IconButton color="inherit" onClick={onThemeChange}>
          {darkMode ? <Brightness7Icon /> : <Brightness4Icon />}
        </IconButton>
      </Toolbar>
    </AppBar>
  );
};
