import React, { useState, useEffect } from 'react';
import { ThemeProvider, createTheme } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import Box from '@mui/material/Box';
import Container from '@mui/material/Container';
import { FileExplorer } from './components/FileExplorer';
import { NodeManager } from './components/NodeManager';
import { SystemStatus } from './components/SystemStatus';
import { AppHeader } from './components/AppHeader';
import { AppTabs } from './components/AppTabs';
import './App.css';

// Create a dark/light theme
const darkTheme = createTheme({
  palette: {
    mode: 'dark',
    primary: {
      main: '#3f51b5',
    },
    secondary: {
      main: '#f50057',
    },
  },
});

const lightTheme = createTheme({
  palette: {
    mode: 'light',
    primary: {
      main: '#3f51b5',
    },
    secondary: {
      main: '#f50057',
    },
  },
});

function App() {
  const [darkMode, setDarkMode] = useState(false);
  const [currentTab, setCurrentTab] = useState(0);

  const handleThemeChange = () => {
    setDarkMode(!darkMode);
  };

  const handleTabChange = (event, newValue) => {
    setCurrentTab(newValue);
  };

  return (
    <ThemeProvider theme={darkMode ? darkTheme : lightTheme}>
      <CssBaseline />
      <Box sx={{ display: 'flex', flexDirection: 'column', minHeight: '100vh' }}>
        <AppHeader darkMode={darkMode} onThemeChange={handleThemeChange} />
        
        <Container maxWidth="lg" sx={{ mt: 4, mb: 4, flexGrow: 1 }}>
          <AppTabs currentTab={currentTab} onTabChange={handleTabChange} />
          
          <Box sx={{ mt: 2 }}>
            {currentTab === 0 && <FileExplorer />}
            {currentTab === 1 && <NodeManager />}
            {currentTab === 2 && <SystemStatus />}
          </Box>
        </Container>
        
        <Box component="footer" sx={{ p: 2, bgcolor: 'background.paper' }}>
          <Container maxWidth="lg">
            <Box sx={{ textAlign: 'center' }}>
              Distributed File System © {new Date().getFullYear()}
            </Box>
          </Container>
        </Box>
      </Box>
    </ThemeProvider>
  );
}

export default App;
