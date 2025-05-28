import React, { useState, useEffect } from 'react';
import axios from 'axios';
import Box from '@mui/material/Box';
import Paper from '@mui/material/Paper';
import Grid from '@mui/material/Grid';
import Typography from '@mui/material/Typography';
import CircularProgress from '@mui/material/CircularProgress';
import Alert from '@mui/material/Alert';
import IconButton from '@mui/material/IconButton';
import RefreshIcon from '@mui/icons-material/Refresh';
import Card from '@mui/material/Card';
import CardContent from '@mui/material/CardContent';
import LinearProgress from '@mui/material/LinearProgress';
import { styled } from '@mui/material/styles';

// Icons
import StorageIcon from '@mui/icons-material/Storage';
import CloudIcon from '@mui/icons-material/Cloud';
import CloudOffIcon from '@mui/icons-material/CloudOff';
import ErrorIcon from '@mui/icons-material/Error';
import SpeedIcon from '@mui/icons-material/Speed';

const API_URL = 'http://localhost:8080/api';

// Styled components
const StatusCard = styled(Card)(({ theme }) => ({
  height: '100%',
  display: 'flex',
  flexDirection: 'column',
  transition: 'transform 0.3s',
  '&:hover': {
    transform: 'translateY(-5px)',
    boxShadow: theme.shadows[6],
  },
}));

const IconWrapper = styled(Box)(({ theme }) => ({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  fontSize: '3rem',
  marginBottom: theme.spacing(2),
}));

export const SystemStatus = () => {
  const [status, setStatus] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  // Fetch system status
  const fetchSystemStatus = async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await axios.get(`${API_URL}/status`);
      setStatus(response.data);
    } catch (err) {
      setError('Failed to fetch system status: ' + (err.response?.data?.error || err.message));
    } finally {
      setLoading(false);
    }
  };

  // Initialize by fetching system status
  useEffect(() => {
    fetchSystemStatus();
    
    // Set up a refresh interval (every 10 seconds)
    const intervalId = setInterval(fetchSystemStatus, 10000);
    
    // Clean up interval on component unmount
    return () => clearInterval(intervalId);
  }, []);

  // Helper to format storage size
  const formatStorage = (bytes) => {
    if (bytes === 0 || !bytes) return '0 Bytes';
    
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  // Helper to calculate storage percentage
  const calculateStoragePercentage = (used, total) => {
    if (!used || !total || total === 0) return 0;
    return Math.round((used / total) * 100);
  };

  return (
    <Box>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
        <Typography variant="h6" component="div">
          System Status
        </Typography>
        <IconButton 
          color="primary" 
          onClick={fetchSystemStatus}
          disabled={loading}
        >
          <RefreshIcon />
        </IconButton>
      </Box>
      
      {loading && !status ? (
        <Box sx={{ display: 'flex', justifyContent: 'center', p: 4 }}>
          <CircularProgress />
        </Box>
      ) : error ? (
        <Alert severity="error">{error}</Alert>
      ) : status ? (
        <Grid container spacing={3}>
          {/* Total Nodes */}
          <Grid item xs={12} sm={6} md={3}>
            <StatusCard>
              <CardContent>
                <IconWrapper>
                  <StorageIcon color="primary" fontSize="inherit" />
                </IconWrapper>
                <Typography variant="h4" component="div" align="center">
                  {status.totalNodes || 0}
                </Typography>
                <Typography variant="body1" color="text.secondary" align="center">
                  Total Nodes
                </Typography>
              </CardContent>
            </StatusCard>
          </Grid>
          
          {/* Active Nodes */}
          <Grid item xs={12} sm={6} md={3}>
            <StatusCard>
              <CardContent>
                <IconWrapper>
                  <CloudIcon sx={{ color: '#4caf50' }} fontSize="inherit" />
                </IconWrapper>
                <Typography variant="h4" component="div" align="center">
                  {status.activeNodes || 0}
                </Typography>
                <Typography variant="body1" color="text.secondary" align="center">
                  Active Nodes
                </Typography>
              </CardContent>
            </StatusCard>
          </Grid>
          
          {/* Inactive Nodes */}
          <Grid item xs={12} sm={6} md={3}>
            <StatusCard>
              <CardContent>
                <IconWrapper>
                  <CloudOffIcon sx={{ color: '#ff9800' }} fontSize="inherit" />
                </IconWrapper>
                <Typography variant="h4" component="div" align="center">
                  {status.inactiveNodes || 0}
                </Typography>
                <Typography variant="body1" color="text.secondary" align="center">
                  Inactive Nodes
                </Typography>
              </CardContent>
            </StatusCard>
          </Grid>
          
          {/* Failed Nodes */}
          <Grid item xs={12} sm={6} md={3}>
            <StatusCard>
              <CardContent>
                <IconWrapper>
                  <ErrorIcon sx={{ color: '#f44336' }} fontSize="inherit" />
                </IconWrapper>
                <Typography variant="h4" component="div" align="center">
                  {status.failedNodes || 0}
                </Typography>
                <Typography variant="body1" color="text.secondary" align="center">
                  Failed Nodes
                </Typography>
              </CardContent>
            </StatusCard>
          </Grid>
          
          {/* Storage Status */}
          <Grid item xs={12}>
            <Paper sx={{ p: 3 }}>
              <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
                <SpeedIcon color="primary" sx={{ mr: 1 }} />
                <Typography variant="h6" component="div">
                  Storage Status
                </Typography>
              </Box>
              
              <Box sx={{ mb: 2 }}>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 0.5 }}>
                  <Typography variant="body2" color="text.secondary">
                    Used: {formatStorage(status.usedStorage || 0)}
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    Available: {formatStorage(status.availableStorage || 0)}
                  </Typography>
                </Box>
                <Box sx={{ width: '100%' }}>
                  <LinearProgress 
                    variant="determinate" 
                    value={calculateStoragePercentage(status.usedStorage, status.totalStorage)} 
                    color={
                      calculateStoragePercentage(status.usedStorage, status.totalStorage) > 90 
                        ? 'error' 
                        : calculateStoragePercentage(status.usedStorage, status.totalStorage) > 70 
                          ? 'warning' 
                          : 'primary'
                    }
                    sx={{ height: 10, borderRadius: 5 }}
                  />
                </Box>
                <Box sx={{ display: 'flex', justifyContent: 'center', mt: 1 }}>
                  <Typography variant="body2" color="text.secondary">
                    {calculateStoragePercentage(status.usedStorage, status.totalStorage)}% of {formatStorage(status.totalStorage || 0)} Used
                  </Typography>
                </Box>
              </Box>
              
              <Typography variant="body2" color="text.secondary" align="center">
                Last updated: {new Date().toLocaleString()}
              </Typography>
            </Paper>
          </Grid>
        </Grid>
      ) : (
        <Alert severity="info">No system status data available</Alert>
      )}
    </Box>
  );
};
