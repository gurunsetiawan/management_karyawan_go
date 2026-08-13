import React, { useState, useEffect, useMemo, useCallback } from 'react';
import {
  Box, Button, IconButton, TextField, AppBar, Toolbar,
  Tooltip, Snackbar, Alert, LinearProgress, InputAdornment,
  Typography, Paper, Grid, Dialog, DialogTitle, DialogContent, DialogActions
} from '@mui/material';
import { DataGrid, GridColDef, GridSortModel, GridToolbar } from '@mui/x-data-grid';
import {
  Add as AddIcon,
  Edit as EditIcon,
  Delete as DeleteIcon,
  Work as WorkIcon,
  Search as SearchIcon,
  Person as PersonIcon,
  Email as EmailIcon,
  Phone as PhoneIcon,
  CloudUpload as CloudUploadIcon
} from '@mui/icons-material';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import api, { Employee } from '../services/api';

// Extend the Employee interface to include required fields
type EmployeeWithRole = Omit<Employee, 'role'> & {
  role: string;
};

interface SnackbarState {
  open: boolean;
  message: string;
  severity: 'success' | 'error' | 'info' | 'warning';
}

const EmployeeList: React.FC = () => {
  const navigate = useNavigate();
  
  // State management
  const [employees, setEmployees] = useState<EmployeeWithRole[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [searchTerm, setSearchTerm] = useState<string>('');
  const [snackbar, setSnackbar] = useState<SnackbarState>({
    open: false,
    message: '',
    severity: 'info'
  });
  const [paginationModel, setPaginationModel] = useState({
    page: 0,
    pageSize: 10,
  });
  const [rowCount, setRowCount] = useState<number>(0);
  const [sortModel, setSortModel] = useState<GridSortModel>([
    { field: 'name', sort: 'asc' },
  ]);
  const [importDialogOpen, setImportDialogOpen] = useState(false);
  const [importFile, setImportFile] = useState<File | null>(null);
  const [importing, setImporting] = useState(false);

  // Fetch employees on component mount
  const fetchEmployees = useCallback(async (searchQuery: string = '', signal?: AbortSignal) => {
    try {
      setLoading(true);
      const res = await api.getEmployees(paginationModel.page + 1, paginationModel.pageSize, searchQuery, { signal });
      const typedEmployees = (res.data || []).map(emp => ({
        ...emp,
        role: emp.role || 'user'
      }));
      setEmployees(typedEmployees);
      setRowCount(res.meta.total);
    } catch (error) {
      if (axios.isCancel(error)) return;
      console.error('Error fetching employees:', error);
      setSnackbar({
        open: true,
        message: 'Failed to fetch employees',
        severity: 'error',
      });
    } finally {
      setLoading(false);
    }
  }, [paginationModel.page, paginationModel.pageSize]);

  useEffect(() => {
    const controller = new AbortController();
    const timer = setTimeout(() => {
      fetchEmployees(searchTerm, controller.signal);
    }, 300);
    return () => {
      clearTimeout(timer);
      controller.abort();
    };
  }, [fetchEmployees, searchTerm]);

  // Event handlers
  const handleDelete = async (id: number) => {
    if (!window.confirm('Apakah Anda yakin ingin menghapus data karyawan ini?')) {
      try {
        await api.deleteEmployee(id);
        await fetchEmployees();
        showSnackbar('Data karyawan berhasil dihapus', 'success');
      } catch (error) {
        showSnackbar('Gagal menghapus data karyawan', 'error');
      }
    }
  };

  const handleEdit = (id: number) => {
    navigate(`/employees/${id}/edit`);
  };

  const handleAddNew = () => {
    navigate('/employees/new');
  };

  const handleDownloadTemplate = () => {
    const csvContent = "Name,Email,Position,Role,Phone,Alamat\nJohn Doe,john@example.com,Software Engineer,admin,0812345678,Jakarta";
    const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
    const link = document.createElement("a");
    const url = URL.createObjectURL(blob);
    link.setAttribute("href", url);
    link.setAttribute("download", "employee_template.csv");
    link.style.visibility = 'hidden';
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files.length > 0) {
      setImportFile(e.target.files[0]);
    }
  };

  const handleImport = async () => {
    if (!importFile) {
      showSnackbar('Silakan pilih file CSV', 'warning');
      return;
    }

    const formData = new FormData();
    formData.append('file', importFile);

    try {
      setImporting(true);
      const token = localStorage.getItem('token');
      const response = await axios.post('http://127.0.0.1:8083/api/employees/import', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
          'Authorization': `Bearer ${token}`
        }
      });
      
      const { success_count, failures } = response.data;
      if (failures && failures.length > 0) {
        showSnackbar(`Berhasil: ${success_count}. Gagal: ${failures.length} baris (Cek console)`, 'warning');
        console.warn("Import Failures:", failures);
      } else {
        showSnackbar(`Berhasil mengimpor ${success_count} karyawan!`, 'success');
      }
      
      setImportDialogOpen(false);
      setImportFile(null);
      fetchEmployees(searchTerm);
    } catch (err: any) {
      showSnackbar(err.response?.data || 'Gagal mengimpor file', 'error');
    } finally {
      setImporting(false);
    }
  };

  const showSnackbar = (message: string, severity: 'success' | 'error' | 'info' | 'warning') => {
    setSnackbar({ open: true, message, severity });
  };

  const handleCloseSnackbar = () => {
    setSnackbar({ ...snackbar, open: false });
  };

  // Table columns
  const columns: GridColDef[] = [
    {
      field: 'name',
      headerName: 'Nama',
      flex: 1,
      renderHeader: () => (
        <Box sx={{ display: 'flex', alignItems: 'center' }}>
          <PersonIcon sx={{ mr: 1 }} />
          <span>Nama</span>
        </Box>
      )
    },
    {
      field: 'position',
      headerName: 'Jabatan',
      flex: 1,
      renderHeader: () => (
        <Box sx={{ display: 'flex', alignItems: 'center' }}>
          <WorkIcon sx={{ mr: 1 }} />
          <span>Jabatan</span>
        </Box>
      )
    },
    {
      field: 'email',
      headerName: 'Email',
      flex: 1,
      renderHeader: () => (
        <Box sx={{ display: 'flex', alignItems: 'center' }}>
          <EmailIcon sx={{ mr: 1 }} />
          <span>Email</span>
        </Box>
      )
    },
    {
      field: 'phone',
      headerName: 'Telepon',
      flex: 1,
      renderHeader: () => (
        <Box sx={{ display: 'flex', alignItems: 'center' }}>
          <PhoneIcon sx={{ mr: 1 }} />
          <span>Telepon</span>
        </Box>
      )
    },
    {
      field: 'actions',
      headerName: 'Aksi',
      sortable: false,
      width: 150,
      renderCell: (params) => (
        <Box sx={{ display: 'flex', gap: 1 }}>
          <Tooltip title="Edit">
            <IconButton
              size="small"
              color="primary"
              onClick={(e) => {
                e.stopPropagation();
                handleEdit((params.row as EmployeeWithRole).id);
              }}
            >
              <EditIcon fontSize="small" />
            </IconButton>
          </Tooltip>
          <Tooltip title="Hapus">
            <IconButton
              size="small"
              color="error"
              onClick={(e) => {
                e.stopPropagation();
                handleDelete((params.row as EmployeeWithRole).id);
              }}
            >
              <DeleteIcon fontSize="small" />
            </IconButton>
          </Tooltip>
        </Box>
      ),
    },
  ];

  return (
    <Box sx={{ height: '100vh', display: 'flex', flexDirection: 'column' }}>
      <AppBar position="static" color="default" elevation={1}>
        <Toolbar>
          <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
            Daftar Karyawan
          </Typography>
        </Toolbar>
      </AppBar>

      <Box sx={{ p: 3, flex: 1, display: 'flex', flexDirection: 'column' }}>
        <Grid container spacing={2} sx={{ mb: 2 }}>
          <Grid size={{ xs: 12, md: 6 }}>
            <TextField
              fullWidth
              variant="outlined"
              placeholder="Cari karyawan..."
              value={searchTerm}
              onChange={(e) => {
                setSearchTerm(e.target.value);
                setPaginationModel(prev => ({ ...prev, page: 0 }));
              }}
              InputProps={{
                startAdornment: (
                  <InputAdornment position="start">
                    <SearchIcon />
                  </InputAdornment>
                ),
              }}
            />
          </Grid>
          <Grid size={{ xs: 12, md: 6 }} sx={{ display: 'flex', justifyContent: 'flex-end', gap: 1 }}>
            <Button
              variant="outlined"
              color="primary"
              onClick={() => setImportDialogOpen(true)}
              startIcon={<CloudUploadIcon />}
            >
              Import CSV
            </Button>
            <Button
              variant="contained"
              color="primary"
              onClick={handleAddNew}
              startIcon={<AddIcon />}
            >
              Tambah Karyawan Baru
            </Button>
          </Grid>
        </Grid>

        <Paper sx={{ p: 2, flexGrow: 1, display: 'flex', flexDirection: 'column' }}>
          <DataGrid
            rows={employees}
            columns={columns}
            rowCount={rowCount}
            paginationMode="server"
            pageSizeOptions={[5, 10, 25]}
            paginationModel={paginationModel}
            onPaginationModelChange={setPaginationModel}
            sortModel={sortModel}
            onSortModelChange={setSortModel}
            loading={loading}
            disableColumnMenu
            disableRowSelectionOnClick
            autoHeight
            slots={{
              toolbar: GridToolbar,
              loadingOverlay: () => <LinearProgress />,
            }}
            slotProps={{
              toolbar: {
                showQuickFilter: true,
              },
            }}
            sx={{
              '& .MuiDataGrid-columnHeaders': {
                backgroundColor: 'primary.light',
                color: 'primary.contrastText',
              },
              '& .MuiDataGrid-cell': {
                borderBottom: '1px solid rgba(224, 224, 224, 1)',
              },
            }}
          />
        </Paper>

        {/* CSV Import Dialog */}
        <Dialog open={importDialogOpen} onClose={() => !importing && setImportDialogOpen(false)}>
          <DialogTitle>Import Data Karyawan</DialogTitle>
          <DialogContent>
            <Typography variant="body2" sx={{ mb: 2 }}>
              Unggah file CSV dengan format kolom: Name, Email, Position, Role, Phone, Alamat.
            </Typography>
            <Button variant="text" onClick={handleDownloadTemplate} sx={{ mb: 2 }}>
              Download Template CSV
            </Button>
            <Box>
              <input
                accept=".csv"
                id="csv-upload"
                type="file"
                onChange={handleFileChange}
                disabled={importing}
              />
            </Box>
            {importing && <LinearProgress sx={{ mt: 2 }} />}
          </DialogContent>
          <DialogActions>
            <Button onClick={() => setImportDialogOpen(false)} disabled={importing}>Batal</Button>
            <Button onClick={handleImport} variant="contained" disabled={!importFile || importing}>
              Upload & Import
            </Button>
          </DialogActions>
        </Dialog>

        {/* Snackbar for notifications */}
        <Snackbar
          open={snackbar.open}
          autoHideDuration={6000}
          onClose={handleCloseSnackbar}
          anchorOrigin={{ vertical: 'top', horizontal: 'center' }}
        >
          <Alert
            onClose={handleCloseSnackbar}
            severity={snackbar.severity}
            sx={{ width: '100%' }}
          >
            {snackbar.message}
          </Alert>
        </Snackbar>
      </Box>
    </Box>
  );
};

export default EmployeeList;
