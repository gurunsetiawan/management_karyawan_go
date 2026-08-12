import React, { useState, useEffect, useMemo, useCallback } from 'react';
import {
  Box, Button, IconButton, TextField, AppBar, Toolbar,
  Tooltip, Snackbar, Alert, LinearProgress, InputAdornment,
  Typography, Paper, Grid
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
  Phone as PhoneIcon
} from '@mui/icons-material';
import { useNavigate } from 'react-router-dom';
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

  // Fetch employees on component mount
  const fetchEmployees = useCallback(async () => {
    try {
      setLoading(true);
      const res = await api.getEmployees(paginationModel.page + 1, paginationModel.pageSize);
      const typedEmployees = (res.data || []).map(emp => ({
        ...emp,
        role: emp.role || 'user'
      }));
      setEmployees(typedEmployees);
      setRowCount(res.meta.total);
    } catch (error) {
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
    fetchEmployees();
  }, [fetchEmployees]);

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

  const showSnackbar = (message: string, severity: 'success' | 'error' | 'info' | 'warning') => {
    setSnackbar({ open: true, message, severity });
  };

  const handleCloseSnackbar = () => {
    setSnackbar({ ...snackbar, open: false });
  };

  // Filter employees based on search term
  const filteredEmployees = useMemo(() => {
    if (!searchTerm) return employees;

    const term = searchTerm.toLowerCase();
    return employees.filter(emp =>
      emp.name.toLowerCase().includes(term) ||
      (emp.position && emp.position.toLowerCase().includes(term)) ||
      (emp.email && emp.email.toLowerCase().includes(term)) ||
      (emp.phone && emp.phone.includes(term))
    );
  }, [employees, searchTerm]);

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
          <Button
            variant="contained"
            color="primary"
            startIcon={<AddIcon />}
            onClick={handleAddNew}
          >
            Tambah Karyawan
          </Button>
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
              onChange={(e) => setSearchTerm(e.target.value)}
              InputProps={{
                startAdornment: (
                  <InputAdornment position="start">
                    <SearchIcon />
                  </InputAdornment>
                ),
              }}
            />
          </Grid>
          <Grid size={{ xs: 12, md: 6 }} sx={{ display: 'flex', justifyContent: 'flex-end' }}>
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
