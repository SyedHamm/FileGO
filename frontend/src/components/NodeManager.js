import React, { useState, useEffect } from 'react';
import axios from 'axios';
import Box from '@mui/material/Box';
import Paper from '@mui/material/Paper';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import Typography from '@mui/material/Typography';
import Button from '@mui/material/Button';
import IconButton from '@mui/material/IconButton';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogContentText from '@mui/material/DialogContentText';
import DialogTitle from '@mui/material/DialogTitle';
import TextField from '@mui/material/TextField';
import MenuItem from '@mui/material/MenuItem';
import CircularProgress from '@mui/material/CircularProgress';
import Alert from '@mui/material/Alert';
import Snackbar from '@mui/material/Snackbar';
import Tooltip from '@mui/material/Tooltip';
import Chip from '@mui/material/Chip';
import LinearProgress from '@mui/material/LinearProgress';

// Icons
import AddIcon from '@mui/icons-material/Add';
import DeleteIcon from '@mui/icons-material/Delete';
import RefreshIcon from '@mui/icons-material/Refresh';
import EditIcon from '@mui/icons-material/Edit';

const API_URL = 'http://localhost:8080/api';

export const NodeManager = () => {
  const [nodes, setNodes] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [registerNodeOpen, setRegisterNodeOpen] = useState(false);
  const [newNode, setNewNode] = useState({ address: '', storageMax: 1024 * 1024 * 1024 }); // 1GB default
  const [editNodeOpen, setEditNodeOpen] = useState(false);
  const [editingNode, setEditingNode] = useState(null);
  const [deleteConfirmOpen, setDeleteConfirmOpen] = useState(false);
  const [nodeToDelete, setNodeToDelete] = useState(null);
  const [notification, setNotification] = useState({ open: false, message: '', severity: 'info' });

  // Fetch nodes
  const fetchNodes = async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await axios.get(`${API_URL}/nodes`);
      setNodes(response.data);
    } catch (err) {
      setError('Failed to fetch nodes: ' + (err.response?.data?.error || err.message));
      setNotification({ 
        open: true, 
        message: 'Failed to fetch nodes', 
        severity: 'error' 
      });
    } finally {
      setLoading(false);
    }
  };

  // Initialize by fetching nodes
  useEffect(() => {
    fetchNodes();
    
    // Set up a refresh interval (every 30 seconds)
    const intervalId = setInterval(fetchNodes, 30000);
    
    // Clean up interval on component unmount
    return () => clearInterval(intervalId);
  }, []);

  // Register node dialog handlers
  const handleRegisterNodeOpen = () => {
    setRegisterNodeOpen(true);
  };

  const handleRegisterNodeClose = () => {
    setRegisterNodeOpen(false);
    setNewNode({ address: '', storageMax: 1024 * 1024 * 1024 });
  };

  const handleRegisterNode = async () => {
    if (!newNode.address || !newNode.storageMax) return;
    
    try {
      await axios.post(`${API_URL}/nodes`, {
        address: newNode.address,
        storageMax: parseInt(newNode.storageMax)
      });
      setNotification({ 
        open: true, 
        message: 'Node registered successfully', 
        severity: 'success' 
      });
      fetchNodes();
    } catch (err) {
      setError('Failed to register node: ' + (err.response?.data?.error || err.message));
      setNotification({ 
        open: true, 
        message: 'Failed to register node', 
        severity: 'error' 
      });
    }
    
    handleRegisterNodeClose();
  };

  // Edit node dialog handlers
  const handleEditNodeOpen = (node) => {
    setEditingNode({
      ...node,
      status: node.status || 'active'
    });
    setEditNodeOpen(true);
  };

  const handleEditNodeClose = () => {
    setEditNodeOpen(false);
    setEditingNode(null);
  };

  const handleEditNode = async () => {
    if (!editingNode) return;
    
    try {
      await axios.put(`${API_URL}/nodes/${editingNode.id}/status`, {
        status: editingNode.status
      });
      setNotification({ 
        open: true, 
        message: 'Node updated successfully', 
        severity: 'success' 
      });
      fetchNodes();
    } catch (err) {
      setError('Failed to update node: ' + (err.response?.data?.error || err.message));
      setNotification({ 
        open: true, 
        message: 'Failed to update node', 
        severity: 'error' 
      });
    }
    
    handleEditNodeClose();
  };

  // Delete node handlers
  const handleDeleteConfirmOpen = (node) => {
    setNodeToDelete(node);
    setDeleteConfirmOpen(true);
  };

  const handleDeleteConfirmClose = () => {
    setDeleteConfirmOpen(false);
    setNodeToDelete(null);
  };

  const handleDeleteNode = async () => {
    if (!nodeToDelete) return;
    
    try {
      await axios.delete(`${API_URL}/nodes/${nodeToDelete.id}`);
      setNotification({ 
        open: true, 
        message: 'Node deleted successfully', 
        severity: 'success' 
      });
      fetchNodes();
    } catch (err) {
      setError('Failed to delete node: ' + (err.response?.data?.error || err.message));
      setNotification({ 
        open: true, 
        message: 'Failed to delete node', 
        severity: 'error' 
      });
    }
    
    handleDeleteConfirmClose();
  };

  // Notification close handler
  const handleNotificationClose = () => {
    setNotification({ ...notification, open: false });
  };

  // Helper to get status chip color
  const getStatusColor = (status) => {
    switch (status) {
      case 'active':
        return 'success';
      case 'inactive':
        return 'warning';
      case 'failed':
        return 'error';
      default:
        return 'default';
    }
  };

  // Helper to format storage size
  const formatStorage = (bytes) => {
    if (bytes === 0) return '0 Bytes';
    
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  // Helper to calculate storage percentage
  const calculateStoragePercentage = (used, max) => {
    return Math.round((used / max) * 100);
  };

  return (
    <Box>
      <Paper sx={{ p: 2, mb: 2 }}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
          <Typography variant="h6" component="div">
            Node Manager
          </Typography>
          <Box>
            <Button 
              variant="contained" 
              color="primary" 
              startIcon={<AddIcon />}
              onClick={handleRegisterNodeOpen}
              sx={{ mr: 1 }}
            >
              Register Node
            </Button>
            <IconButton 
              color="primary" 
              onClick={fetchNodes}
              disabled={loading}
            >
              <RefreshIcon />
            </IconButton>
          </Box>
        </Box>
        
        {loading ? (
          <Box sx={{ display: 'flex', justifyContent: 'center', p: 4 }}>
            <CircularProgress />
          </Box>
        ) : error ? (
          <Alert severity="error">{error}</Alert>
        ) : nodes.length === 0 ? (
          <Box sx={{ p: 4, textAlign: 'center' }}>
            <Typography variant="body1" color="text.secondary">
              No nodes registered
            </Typography>
          </Box>
        ) : (
          <TableContainer>
            <Table sx={{ minWidth: 650 }} aria-label="nodes table">
              <TableHead>
                <TableRow>
                  <TableCell>ID</TableCell>
                  <TableCell>Address</TableCell>
                  <TableCell>Status</TableCell>
                  <TableCell>Storage</TableCell>
                  <TableCell>Last Seen</TableCell>
                  <TableCell>Actions</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {nodes.map((node) => (
                  <TableRow key={node.id}>
                    <TableCell component="th" scope="row">
                      {node.id.substring(0, 8)}...
                    </TableCell>
                    <TableCell>{node.address}</TableCell>
                    <TableCell>
                      <Chip 
                        label={node.status || 'active'} 
                        color={getStatusColor(node.status || 'active')} 
                        size="small" 
                      />
                    </TableCell>
                    <TableCell>
                      <Box sx={{ display: 'flex', alignItems: 'center', flexDirection: 'column' }}>
                        <Box sx={{ width: '100%', mr: 1 }}>
                          <LinearProgress 
                            variant="determinate" 
                            value={calculateStoragePercentage(node.storageUsed, node.storageMax)} 
                            color={
                              calculateStoragePercentage(node.storageUsed, node.storageMax) > 90 
                                ? 'error' 
                                : calculateStoragePercentage(node.storageUsed, node.storageMax) > 70 
                                  ? 'warning' 
                                  : 'primary'
                            }
                          />
                        </Box>
                        <Box sx={{ minWidth: 35, mt: 1 }}>
                          <Typography variant="body2" color="text.secondary">
                            {formatStorage(node.storageUsed)} / {formatStorage(node.storageMax)}
                          </Typography>
                        </Box>
                      </Box>
                    </TableCell>
                    <TableCell>
                      {new Date(node.lastSeen).toLocaleString()}
                    </TableCell>
                    <TableCell>
                      <Tooltip title="Edit Node">
                        <IconButton 
                          color="primary" 
                          onClick={() => handleEditNodeOpen(node)}
                          size="small"
                          sx={{ mr: 1 }}
                        >
                          <EditIcon fontSize="small" />
                        </IconButton>
                      </Tooltip>
                      <Tooltip title="Delete Node">
                        <IconButton 
                          color="error" 
                          onClick={() => handleDeleteConfirmOpen(node)}
                          size="small"
                        >
                          <DeleteIcon fontSize="small" />
                        </IconButton>
                      </Tooltip>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </TableContainer>
        )}
      </Paper>

      {/* Register Node Dialog */}
      <Dialog open={registerNodeOpen} onClose={handleRegisterNodeClose}>
        <DialogTitle>Register New Node</DialogTitle>
        <DialogContent>
          <DialogContentText>
            Enter the details for the new node:
          </DialogContentText>
          <TextField
            autoFocus
            margin="dense"
            id="address"
            label="Node Address"
            type="text"
            fullWidth
            variant="outlined"
            value={newNode.address}
            onChange={(e) => setNewNode({ ...newNode, address: e.target.value })}
            sx={{ mb: 2 }}
            placeholder="http://nodeaddress:port"
          />
          <TextField
            margin="dense"
            id="storageMax"
            label="Max Storage (bytes)"
            type="number"
            fullWidth
            variant="outlined"
            value={newNode.storageMax}
            onChange={(e) => setNewNode({ ...newNode, storageMax: e.target.value })}
            inputProps={{ min: 1 }}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={handleRegisterNodeClose} color="primary">
            Cancel
          </Button>
          <Button 
            onClick={handleRegisterNode} 
            color="primary"
            disabled={!newNode.address || !newNode.storageMax}
          >
            Register
          </Button>
        </DialogActions>
      </Dialog>

      {/* Edit Node Dialog */}
      <Dialog open={editNodeOpen} onClose={handleEditNodeClose}>
        <DialogTitle>Edit Node</DialogTitle>
        <DialogContent>
          <DialogContentText>
            Update node status:
          </DialogContentText>
          <TextField
            select
            margin="dense"
            id="status"
            label="Status"
            fullWidth
            variant="outlined"
            value={editingNode?.status || 'active'}
            onChange={(e) => setEditingNode({ ...editingNode, status: e.target.value })}
            sx={{ mb: 2 }}
          >
            <MenuItem value="active">Active</MenuItem>
            <MenuItem value="inactive">Inactive</MenuItem>
            <MenuItem value="failed">Failed</MenuItem>
          </TextField>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleEditNodeClose} color="primary">
            Cancel
          </Button>
          <Button onClick={handleEditNode} color="primary">
            Update
          </Button>
        </DialogActions>
      </Dialog>

      {/* Delete Confirmation Dialog */}
      <Dialog open={deleteConfirmOpen} onClose={handleDeleteConfirmClose}>
        <DialogTitle>Confirm Deletion</DialogTitle>
        <DialogContent>
          <DialogContentText>
            Are you sure you want to delete this node? This action cannot be undone.
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleDeleteConfirmClose} color="primary">
            Cancel
          </Button>
          <Button onClick={handleDeleteNode} color="error">
            Delete
          </Button>
        </DialogActions>
      </Dialog>

      {/* Notification Snackbar */}
      <Snackbar
        open={notification.open}
        autoHideDuration={6000}
        onClose={handleNotificationClose}
        anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
      >
        <Alert 
          onClose={handleNotificationClose} 
          severity={notification.severity} 
          sx={{ width: '100%' }}
        >
          {notification.message}
        </Alert>
      </Snackbar>
    </Box>
  );
};
