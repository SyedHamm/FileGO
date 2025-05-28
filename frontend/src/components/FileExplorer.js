import React, { useState, useEffect } from 'react';
import axios from 'axios';
import Box from '@mui/material/Box';
import Paper from '@mui/material/Paper';
import List from '@mui/material/List';
import ListItem from '@mui/material/ListItem';
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemText from '@mui/material/ListItemText';
import Typography from '@mui/material/Typography';
import Breadcrumbs from '@mui/material/Breadcrumbs';
import Link from '@mui/material/Link';
import IconButton from '@mui/material/IconButton';
import Button from '@mui/material/Button';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogContentText from '@mui/material/DialogContentText';
import DialogTitle from '@mui/material/DialogTitle';
import TextField from '@mui/material/TextField';
import MenuItem from '@mui/material/MenuItem';
import Menu from '@mui/material/Menu';
import Alert from '@mui/material/Alert';
import Snackbar from '@mui/material/Snackbar';
import CircularProgress from '@mui/material/CircularProgress';
import Divider from '@mui/material/Divider';
import Tooltip from '@mui/material/Tooltip';

// Icons
import FolderIcon from '@mui/icons-material/Folder';
import InsertDriveFileIcon from '@mui/icons-material/InsertDriveFile';
import CreateNewFolderIcon from '@mui/icons-material/CreateNewFolder';
import UploadFileIcon from '@mui/icons-material/UploadFile';
import HomeIcon from '@mui/icons-material/Home';
import MoreVertIcon from '@mui/icons-material/MoreVert';
import DeleteIcon from '@mui/icons-material/Delete';
import FileCopyIcon from '@mui/icons-material/FileCopy';
import GetAppIcon from '@mui/icons-material/GetApp';
import AutorenewIcon from '@mui/icons-material/Autorenew';

import { useDropzone } from 'react-dropzone';

const API_URL = 'http://localhost:8080/api';

export const FileExplorer = () => {
  const [currentPath, setCurrentPath] = useState('');
  const [files, setFiles] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [createFolderOpen, setCreateFolderOpen] = useState(false);
  const [newFolderName, setNewFolderName] = useState('');
  const [uploadOpen, setUploadOpen] = useState(false);
  const [contextMenu, setContextMenu] = useState(null);
  const [selectedFile, setSelectedFile] = useState(null);
  const [replicateOpen, setReplicateOpen] = useState(false);
  const [replicaCount, setReplicaCount] = useState(1);
  const [notification, setNotification] = useState({ open: false, message: '', severity: 'info' });

  // Fetch files from the current path
  const fetchFiles = async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await axios.get(`${API_URL}/files?path=${currentPath}`);
      setFiles(response.data || []);
    } catch (err) {
      setError('Failed to fetch files: ' + (err.response?.data?.error || err.message));
      setNotification({ 
        open: true, 
        message: 'Failed to fetch files', 
        severity: 'error' 
      });
      // Initialize with empty array on error
      setFiles([]);
    } finally {
      setLoading(false);
    }
  };

  // Initialize by fetching root files
  useEffect(() => {
    fetchFiles();
  }, [currentPath]);

  // Navigate to a directory
  const navigateToDirectory = (path) => {
    setCurrentPath(path);
  };

  // Handle breadcrumb navigation
  const handleBreadcrumbClick = (path) => {
    navigateToDirectory(path);
  };

  // Create folder dialog handlers
  const handleCreateFolderOpen = () => {
    setCreateFolderOpen(true);
  };

  const handleCreateFolderClose = () => {
    setCreateFolderOpen(false);
    setNewFolderName('');
  };

  const handleCreateFolder = async () => {
    if (!newFolderName) return;
    
    try {
      await axios.post(`${API_URL}/directories/${currentPath ? currentPath + '/' : ''}${newFolderName}`);
      setNotification({ 
        open: true, 
        message: 'Folder created successfully', 
        severity: 'success' 
      });
      fetchFiles();
    } catch (err) {
      setError('Failed to create folder: ' + (err.response?.data?.error || err.message));
      setNotification({ 
        open: true, 
        message: 'Failed to create folder', 
        severity: 'error' 
      });
    }
    
    handleCreateFolderClose();
  };

  // Upload dialog handlers
  const handleUploadOpen = () => {
    setUploadOpen(true);
  };

  const handleUploadClose = () => {
    setUploadOpen(false);
  };

  // File dropzone configuration
  const { getRootProps, getInputProps, acceptedFiles } = useDropzone({
    accept: {
      'application/octet-stream': [],
      'text/plain': ['.txt'],
      'application/pdf': ['.pdf'],
      'image/*': ['.jpeg', '.jpg', '.png', '.gif'],
      'application/msword': ['.doc', '.docx'],
      'application/vnd.ms-excel': ['.xls', '.xlsx'],
    },
    multiple: true
  });

  const handleFileUpload = async () => {
    if (acceptedFiles.length === 0) return;
    
    try {
      for (const file of acceptedFiles) {
        const formData = new FormData();
        formData.append('file', file);
        
        await axios.post(
          `${API_URL}/files/${currentPath ? currentPath + '/' : ''}${file.name}`, 
          formData, 
          {
            headers: {
              'Content-Type': 'multipart/form-data'
            }
          }
        );
      }
      
      setNotification({ 
        open: true, 
        message: 'Files uploaded successfully', 
        severity: 'success' 
      });
      fetchFiles();
    } catch (err) {
      setError('Failed to upload files: ' + (err.response?.data?.error || err.message));
      setNotification({ 
        open: true, 
        message: 'Failed to upload files', 
        severity: 'error' 
      });
    }
    
    handleUploadClose();
  };

  // Context menu handlers
  const handleContextMenu = (event, file) => {
    event.preventDefault();
    setContextMenu({ mouseX: event.clientX - 2, mouseY: event.clientY - 4 });
    setSelectedFile(file);
  };

  const handleContextMenuClose = () => {
    setContextMenu(null);
    setSelectedFile(null);
  };

  // Delete file handler
  const handleDeleteFile = async () => {
    if (!selectedFile) return;
    
    try {
      await axios.delete(`${API_URL}/files/${selectedFile.path}`);
      setNotification({ 
        open: true, 
        message: 'File deleted successfully', 
        severity: 'success' 
      });
      fetchFiles();
    } catch (err) {
      setError('Failed to delete file: ' + (err.response?.data?.error || err.message));
      setNotification({ 
        open: true, 
        message: 'Failed to delete file', 
        severity: 'error' 
      });
    }
    
    handleContextMenuClose();
  };

  // Download file handler
  const handleDownloadFile = () => {
    if (!selectedFile) return;
    
    window.open(`${API_URL}/files/${selectedFile.path}?download=true`, '_blank');
    handleContextMenuClose();
  };

  // Replicate file dialog handlers
  const handleReplicateOpen = () => {
    setReplicateOpen(true);
  };

  const handleReplicateClose = () => {
    setReplicateOpen(false);
    setReplicaCount(1);
  };

  const handleReplicateFile = async () => {
    if (!selectedFile || replicaCount < 1) return;
    
    try {
      await axios.put(`${API_URL}/replicate/${selectedFile.path}?replicas=${replicaCount}`);
      setNotification({ 
        open: true, 
        message: 'Replication factor set successfully', 
        severity: 'success' 
      });
    } catch (err) {
      setError('Failed to set replication factor: ' + (err.response?.data?.error || err.message));
      setNotification({ 
        open: true, 
        message: 'Failed to set replication factor', 
        severity: 'error' 
      });
    }
    
    handleReplicateClose();
    handleContextMenuClose();
  };

  // Notification close handler
  const handleNotificationClose = () => {
    setNotification({ ...notification, open: false });
  };

  // Generate breadcrumbs for navigation
  const renderBreadcrumbs = () => {
    const paths = currentPath.split('/').filter(Boolean);
    const breadcrumbs = [];
    
    breadcrumbs.push(
      <Link
        key="home"
        underline="hover"
        color="inherit"
        onClick={() => handleBreadcrumbClick('')}
        className="breadcrumb-item"
      >
        <HomeIcon sx={{ mr: 0.5 }} fontSize="inherit" />
        Home
      </Link>
    );
    
    let currentPathAccumulator = '';
    
    for (let i = 0; i < paths.length; i++) {
      currentPathAccumulator += (i === 0 ? '' : '/') + paths[i];
      
      breadcrumbs.push(
        <Link
          key={paths[i]}
          underline="hover"
          color="inherit"
          onClick={() => handleBreadcrumbClick(currentPathAccumulator)}
          className="breadcrumb-item"
        >
          {paths[i]}
        </Link>
      );
    }
    
    return (
      <Breadcrumbs separator="›" aria-label="breadcrumb">
        {breadcrumbs}
      </Breadcrumbs>
    );
  };

  return (
    <Box>
      <Paper sx={{ p: 2, mb: 2 }}>
        <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
          {renderBreadcrumbs()}
        </Box>
        
        <Box sx={{ display: 'flex', justifyContent: 'flex-start', mb: 2 }}>
          <Button 
            variant="contained" 
            color="primary" 
            startIcon={<CreateNewFolderIcon />}
            onClick={handleCreateFolderOpen}
            sx={{ mr: 1 }}
          >
            New Folder
          </Button>
          <Button 
            variant="contained" 
            color="secondary" 
            startIcon={<UploadFileIcon />}
            onClick={handleUploadOpen}
          >
            Upload Files
          </Button>
        </Box>
        
        <Divider sx={{ mb: 2 }} />
        
        {loading ? (
          <Box sx={{ display: 'flex', justifyContent: 'center', p: 4 }}>
            <CircularProgress />
          </Box>
        ) : error ? (
          <Alert severity="error">{error}</Alert>
        ) : !files || files.length === 0 ? (
          <Box sx={{ p: 4, textAlign: 'center' }}>
            <Typography variant="body1" color="text.secondary">
              This folder is empty
            </Typography>
          </Box>
        ) : (
          <List>
            {(files || []).map((file) => (
              <ListItem
                key={file.path}
                className="file-list-item"
                onClick={file.isDir ? () => navigateToDirectory(file.path) : undefined}
                onContextMenu={(e) => handleContextMenu(e, file)}
                secondaryAction={
                  <IconButton edge="end" onClick={(e) => {
                    e.stopPropagation();
                    handleContextMenu(e, file);
                  }}>
                    <MoreVertIcon />
                  </IconButton>
                }
              >
                <ListItemIcon>
                  {file.isDir ? <FolderIcon color="primary" /> : <InsertDriveFileIcon color="action" />}
                </ListItemIcon>
                <ListItemText
                  primary={file.name}
                  secondary={
                    <React.Fragment>
                      <Typography
                        component="span"
                        variant="body2"
                        color="text.primary"
                      >
                        {file.isDir ? 'Folder' : `File • ${formatFileSize(file.size)}`}
                      </Typography>
                      {file.replicas > 1 && ` • ${file.replicas} replicas`}
                    </React.Fragment>
                  }
                />
              </ListItem>
            ))}
          </List>
        )}
      </Paper>

      {/* Create Folder Dialog */}
      <Dialog open={createFolderOpen} onClose={handleCreateFolderClose}>
        <DialogTitle>Create New Folder</DialogTitle>
        <DialogContent>
          <DialogContentText>
            Enter a name for the new folder:
          </DialogContentText>
          <TextField
            autoFocus
            margin="dense"
            id="name"
            label="Folder Name"
            type="text"
            fullWidth
            variant="outlined"
            value={newFolderName}
            onChange={(e) => setNewFolderName(e.target.value)}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCreateFolderClose} color="primary">
            Cancel
          </Button>
          <Button onClick={handleCreateFolder} color="primary" disabled={!newFolderName}>
            Create
          </Button>
        </DialogActions>
      </Dialog>

      {/* Upload Files Dialog */}
      <Dialog open={uploadOpen} onClose={handleUploadClose} maxWidth="sm" fullWidth>
        <DialogTitle>Upload Files</DialogTitle>
        <DialogContent>
          <DialogContentText>
            Drag and drop files here or click to select files:
          </DialogContentText>
          <Box 
            {...getRootProps()} 
            className="upload-dropzone"
            sx={{ mt: 2 }}
          >
            <input {...getInputProps()} />
            <Typography>Drag and drop files here, or click to select files</Typography>
          </Box>
          {acceptedFiles.length > 0 && (
            <Box sx={{ mt: 2 }}>
              <Typography variant="subtitle1">Selected Files:</Typography>
              <List dense>
                {acceptedFiles.map((file, index) => (
                  <ListItem key={index}>
                    <ListItemIcon>
                      <InsertDriveFileIcon />
                    </ListItemIcon>
                    <ListItemText 
                      primary={file.name} 
                      secondary={formatFileSize(file.size)} 
                    />
                  </ListItem>
                ))}
              </List>
            </Box>
          )}
        </DialogContent>
        <DialogActions>
          <Button onClick={handleUploadClose} color="primary">
            Cancel
          </Button>
          <Button 
            onClick={handleFileUpload} 
            color="primary"
            disabled={acceptedFiles.length === 0}
          >
            Upload
          </Button>
        </DialogActions>
      </Dialog>

      {/* Context Menu */}
      <Menu
        open={contextMenu !== null}
        onClose={handleContextMenuClose}
        anchorReference="anchorPosition"
        anchorPosition={
          contextMenu !== null
            ? { top: contextMenu.mouseY, left: contextMenu.mouseX }
            : undefined
        }
      >
        {selectedFile && !selectedFile.isDir && (
          <MenuItem onClick={handleDownloadFile}>
            <ListItemIcon>
              <GetAppIcon fontSize="small" />
            </ListItemIcon>
            <ListItemText>Download</ListItemText>
          </MenuItem>
        )}
        <MenuItem onClick={handleDeleteFile}>
          <ListItemIcon>
            <DeleteIcon fontSize="small" />
          </ListItemIcon>
          <ListItemText>Delete</ListItemText>
        </MenuItem>
        {selectedFile && !selectedFile.isDir && (
          <MenuItem onClick={handleReplicateOpen}>
            <ListItemIcon>
              <AutorenewIcon fontSize="small" />
            </ListItemIcon>
            <ListItemText>Set Replication Factor</ListItemText>
          </MenuItem>
        )}
      </Menu>

      {/* Replication Dialog */}
      <Dialog open={replicateOpen} onClose={handleReplicateClose}>
        <DialogTitle>Set Replication Factor</DialogTitle>
        <DialogContent>
          <DialogContentText>
            Set the number of replicas for this file:
          </DialogContentText>
          <TextField
            autoFocus
            margin="dense"
            id="replicas"
            label="Number of Replicas"
            type="number"
            fullWidth
            variant="outlined"
            value={replicaCount}
            onChange={(e) => setReplicaCount(Math.max(1, parseInt(e.target.value) || 1))}
            inputProps={{ min: 1 }}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={handleReplicateClose} color="primary">
            Cancel
          </Button>
          <Button onClick={handleReplicateFile} color="primary">
            Set
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

// Helper function to format file size
const formatFileSize = (bytes) => {
  if (bytes === 0) return '0 Bytes';
  
  const k = 1024;
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
};
