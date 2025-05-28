import React from 'react';
import Tabs from '@mui/material/Tabs';
import Tab from '@mui/material/Tab';
import Box from '@mui/material/Box';
import FolderIcon from '@mui/icons-material/Folder';
import CloudIcon from '@mui/icons-material/Cloud';
import DashboardIcon from '@mui/icons-material/Dashboard';

export const AppTabs = ({ currentTab, onTabChange }) => {
  return (
    <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
      <Tabs
        value={currentTab}
        onChange={onTabChange}
        aria-label="app navigation tabs"
        centered
      >
        <Tab 
          icon={<FolderIcon />} 
          label="File Explorer" 
          iconPosition="start"
        />
        <Tab 
          icon={<CloudIcon />} 
          label="Node Manager" 
          iconPosition="start"
        />
        <Tab 
          icon={<DashboardIcon />} 
          label="System Status" 
          iconPosition="start"
        />
      </Tabs>
    </Box>
  );
};
