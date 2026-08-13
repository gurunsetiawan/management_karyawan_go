import React, { useState, useEffect, useCallback } from 'react';
import {
  Box, Button, IconButton, TextField, Tooltip, Snackbar, Alert,
  Typography, Card, Avatar, Chip, ToggleButton, ToggleButtonGroup,
  Dialog, DialogTitle, DialogContent, DialogActions, InputAdornment,
  CircularProgress
} from '@mui/material';
import { DataGrid, GridColDef, GridSortModel } from '@mui/x-data-grid';
import {
  Add as AddIcon,
  Edit as EditIcon,
  Delete as DeleteIcon,
  Search as SearchIcon,
  Refresh as RefreshIcon,
  CloudUpload as CloudUploadIcon,
  Print as PrintIcon,
  Download as DownloadIcon,
  FilterList as FilterListIcon,
  ViewColumn as ViewColumnIcon
} from '@mui/icons-material';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import api, { Employee } from '../services/api';

type EmployeeWithRole = Omit<Employee, 'role'> & {
  role: string;
};

interface SnackbarState {
  open: boolean;
  message: string;
  severity: 'success' | 'error' | 'info' | 'warning';
}

function stringToColor(string: string) {
  let hash = 0;
  for (let i = 0; i < string.length; i += 1) {
    hash = string.charCodeAt(i) + ((hash << 5) - hash);
  }
  let color = '#';
  for (let i = 0; i < 3; i += 1) {
    const value = (hash >> (i * 8)) & 0xff;
    color += `00${value.toString(16)}`.slice(-2);
  }
  return color;
}

function stringAvatar(name: string) {
  const parts = name.split(' ');
  const init = parts.length > 1 ? `${parts[0][0]}${parts[1][0]}` : name[0];
  return {
    sx: {
      bgcolor: stringToColor(name),
      width: 32,
      height: 32,
      fontSize: '0.8rem',
      fontWeight: 600,
    },
    children: init ? init.toUpperCase() : 'E',
  };
}

const EmployeeList: React.FC = () => {
  const navigate = useNavigate();

  const [employees, setEmployees] = useState<EmployeeWithRole[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [searchTerm, setSearchTerm] = useState<string>('');
  const [statusFilter, setStatusFilter] = useState<string>('all');
  const [snackbar, setSnackbar] = useState<SnackbarState>({
    open: false,
    message: '',
    severity: 'info',
  });
  const [paginationModel, setPaginationModel] = useState({
    page: 0,
    pageSize: 10,
  });
  const [rowCount, setRowCount] = useState<number>(0);
  const [importDialogOpen, setImportDialogOpen] = useState(false);
  const [importFile, setImportFile] = useState<File | null>(null);
  const [importing, setImporting] = useState(false);

  const fetchEmployees = useCallback(async (searchQuery: string = '', status: string = 'all', signal?: AbortSignal) => {
    try {
      setLoading(true);
      const res = await api.getEmployees(paginationModel.page + 1, paginationModel.pageSize, searchQuery, status, { signal });
      const typedEmployees = (res.data || []).map(emp => ({
        ...emp,
        role: emp.role || 'karyawan',
      }));
      setEmployees(typedEmployees);
      setRowCount(res.meta.total);
    } catch (error) {
      if (axios.isCancel(error)) return;
      console.error('Error fetching employees:', error);
      showSnackbar('Gagal mengambil data karyawan', 'error');
    } finally {
      setLoading(false);
    }
  }, [paginationModel.page, paginationModel.pageSize]);

  useEffect(() => {
    const controller = new AbortController();
    const timer = setTimeout(() => {
      fetchEmployees(searchTerm, statusFilter, controller.signal);
    }, 300);
    return () => {
      clearTimeout(timer);
      controller.abort();
    };
  }, [fetchEmployees, searchTerm, statusFilter]);

  const handleDelete = async (id: number) => {
    if (window.confirm('Apakah Anda yakin ingin menghapus data karyawan ini?')) {
      try {
        await api.deleteEmployee(id);
        await fetchEmployees(searchTerm, statusFilter);
        showSnackbar('Data karyawan berhasil dihapus', 'success');
      } catch (error) {
        showSnackbar('Gagal menghapus data karyawan', 'error');
      }
    }
  };

  const handleEdit = (id: number) => {
    navigate(`/employees/${id}/edit`);
  };

  const handleExportCSV = async () => {
    try {
      const token = localStorage.getItem('token');
      const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8083';
      const response = await axios.get(`${API_URL}/api/employees/export`, {
        headers: { Authorization: `Bearer ${token}` },
        responseType: 'blob',
      });
      const url = window.URL.createObjectURL(new Blob([response.data]));
      const link = document.createElement('a');
      link.href = url;
      link.setAttribute('download', 'employees.csv');
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      window.URL.revokeObjectURL(url);
    } catch (err) {
      showSnackbar('Gagal mengekspor data', 'error');
    }
  };

  const handlePrint = () => {
    window.print();
  };

  const handleAddNew = () => {
    navigate('/employees/new');
  };

  const handleDownloadTemplate = () => {
    const csvContent = 'Name,Email,Position,Role,Phone,Alamat\nJohn Doe,john@example.com,Software Engineer,admin,0812345678,Jakarta';
    const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
    const link = document.createElement('a');
    const url = URL.createObjectURL(blob);
    link.setAttribute('href', url);
    link.setAttribute('download', 'employee_template.csv');
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
      const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8083';
      const response = await axios.post(`${API_URL}/api/employees/import`, formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
          Authorization: `Bearer ${token}`,
        },
      });

      const { success_count, failures } = response.data;
      if (failures && failures.length > 0) {
        showSnackbar(`Berhasil: ${success_count}. Gagal: ${failures.length} baris`, 'warning');
      } else {
        showSnackbar(`Berhasil mengimpor ${success_count} karyawan!`, 'success');
      }

      setImportDialogOpen(false);
      setImportFile(null);
      fetchEmployees(searchTerm, statusFilter);
    } catch (err: any) {
      const errorMsg = typeof err.response?.data === 'string' ? err.response.data : 'Gagal mengimpor file';
      showSnackbar(errorMsg, 'error');
    } finally {
      setImporting(false);
    }
  };

  const showSnackbar = (message: string, severity: 'success' | 'error' | 'info' | 'warning') => {
    setSnackbar({ open: true, message, severity });
  };

  const columns: GridColDef[] = [
    {
      field: 'id',
      headerName: 'ID',
      width: 70,
      align: 'center',
      headerAlign: 'center',
    },
    {
      field: 'name',
      headerName: 'Name',
      flex: 1.5,
      minWidth: 180,
      renderCell: (params) => (
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1.5, height: '100%' }}>
          <Avatar {...stringAvatar(params.row.name)} />
          <Typography variant="body2" sx={{ fontWeight: 500 }}>
            {params.row.name}
          </Typography>
        </Box>
      ),
    },
    {
      field: 'email',
      headerName: 'Email',
      flex: 1.5,
      minWidth: 200,
    },
    {
      field: 'position',
      headerName: 'Department',
      flex: 1.2,
      minWidth: 150,
      renderCell: (params) => (
        <Typography variant="body2" sx={{ color: 'text.secondary' }}>
          {params.value || '-'}
        </Typography>
      ),
    },
    {
      field: 'role',
      headerName: 'Role',
      flex: 1,
      minWidth: 120,
      renderCell: (params) => {
        const role = (params.value || '').toLowerCase();
        let color: 'error' | 'warning' | 'default' = 'default';
        if (role === 'admin') color = 'error';
        if (role === 'manager') color = 'warning';
        return (
          <Chip
            label={params.value || 'Karyawan'}
            size="small"
            color={color}
            variant="outlined"
            sx={{ textTransform: 'capitalize', fontWeight: 500 }}
          />
        );
      },
    },
    {
      field: 'phone',
      headerName: 'Phone',
      flex: 1.2,
      minWidth: 140,
    },
    {
      field: 'alamat',
      headerName: 'Alamat',
      flex: 1.5,
      minWidth: 180,
    },
    {
      field: 'actions',
      headerName: 'Actions',
      width: 100,
      sortable: false,
      filterable: false,
      align: 'center',
      headerAlign: 'center',
      renderCell: (params) => (
        <Box sx={{ display: 'flex', gap: 0.5, justifyContent: 'center', alignItems: 'center', height: '100%' }}>
          <Tooltip title="Edit">
            <IconButton size="small" onClick={() => handleEdit(params.row.id)} color="default">
              <EditIcon fontSize="small" />
            </IconButton>
          </Tooltip>
          <Tooltip title="Hapus">
            <IconButton size="small" onClick={() => handleDelete(params.row.id)} color="default">
              <DeleteIcon fontSize="small" />
            </IconButton>
          </Tooltip>
        </Box>
      ),
    },
  ];

  return (
    <Box>
      {/* Top Header Section matching MUI Template */}
      <Box sx={{ mb: 3, display: 'flex', alignItems: 'flex-start', justifyContent: 'space-between', flexWrap: 'wrap', gap: 2 }}>
        <Box>
          <Typography variant="caption" sx={{ color: 'text.secondary', fontWeight: 500, display: 'block', mb: 0.5 }}>
            Employees
          </Typography>
          <Typography variant="h4" component="h1">
            Employees
          </Typography>
        </Box>

        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
          <Tooltip title="Refresh">
            <IconButton
              onClick={() => fetchEmployees(searchTerm, statusFilter)}
              sx={{
                border: '1px solid',
                borderColor: 'divider',
                borderRadius: '8px',
                p: 1,
              }}
            >
              <RefreshIcon fontSize="small" />
            </IconButton>
          </Tooltip>

          <Button
            variant="contained"
            color="primary"
            onClick={handleAddNew}
            startIcon={<AddIcon />}
            sx={{
              borderRadius: '8px',
              px: 2.5,
              py: 0.9,
              fontWeight: 600,
            }}
          >
            + Create
          </Button>
        </Box>
      </Box>

      {/* Main Card Container */}
      <Card sx={{ p: 0, overflow: 'hidden' }}>
        {/* Card Toolbar Top */}
        <Box
          sx={{
            p: 2,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'space-between',
            flexWrap: 'wrap',
            gap: 2,
            borderBottom: '1px solid',
            borderColor: 'divider',
          }}
        >
          {/* Search Bar & Filter Pills */}
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, flexWrap: 'wrap', flexGrow: 1 }}>
            <TextField
              size="small"
              placeholder="Cari karyawan..."
              value={searchTerm}
              onChange={(e) => {
                setSearchTerm(e.target.value);
                setPaginationModel((prev) => ({ ...prev, page: 0 }));
              }}
              InputProps={{
                startAdornment: (
                  <InputAdornment position="start">
                    <SearchIcon fontSize="small" sx={{ color: 'text.secondary' }} />
                  </InputAdornment>
                ),
              }}
              sx={{ width: { xs: '100%', sm: 260 } }}
            />

            <ToggleButtonGroup
              value={statusFilter}
              exclusive
              onChange={(e, newStatus) => {
                if (newStatus !== null) {
                  setStatusFilter(newStatus);
                  setPaginationModel((prev) => ({ ...prev, page: 0 }));
                }
              }}
              size="small"
              sx={{
                '& .MuiToggleButton-root': {
                  px: 2,
                  py: 0.5,
                  fontSize: '0.8125rem',
                  textTransform: 'none',
                  fontWeight: 500,
                },
              }}
            >
              <ToggleButton value="all">Semua</ToggleButton>
              <ToggleButton value="active">Aktif</ToggleButton>
              <ToggleButton value="inactive">Non-Aktif</ToggleButton>
            </ToggleButtonGroup>
          </Box>

          {/* Action Icons Right Inside Card Header */}
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
            <Tooltip title="Print Table">
              <IconButton size="small" onClick={handlePrint} sx={{ border: '1px solid', borderColor: 'divider', borderRadius: '6px' }}>
                <PrintIcon fontSize="small" />
              </IconButton>
            </Tooltip>

            <Tooltip title="Export CSV">
              <IconButton size="small" onClick={handleExportCSV} sx={{ border: '1px solid', borderColor: 'divider', borderRadius: '6px' }}>
                <DownloadIcon fontSize="small" />
              </IconButton>
            </Tooltip>

            <Tooltip title="Import CSV">
              <IconButton size="small" onClick={() => setImportDialogOpen(true)} sx={{ border: '1px solid', borderColor: 'divider', borderRadius: '6px' }}>
                <CloudUploadIcon fontSize="small" />
              </IconButton>
            </Tooltip>
          </Box>
        </Box>

        {/* DataGrid Component */}
        <Box sx={{ height: 550, width: '100%' }}>
          <DataGrid
            rows={employees}
            columns={columns}
            loading={loading}
            rowCount={rowCount}
            paginationMode="server"
            paginationModel={paginationModel}
            onPaginationModelChange={setPaginationModel}
            pageSizeOptions={[5, 10, 25, 50]}
            disableRowSelectionOnClick
            sx={{
              border: 'none',
              '& .MuiDataGrid-cell': {
                borderBottom: '1px solid',
                borderColor: 'divider',
              },
              '& .MuiDataGrid-columnHeaders': {
                borderBottom: '1px solid',
                borderColor: 'divider',
              },
              '& .MuiDataGrid-footerContainer': {
                borderTop: '1px solid',
                borderColor: 'divider',
              },
            }}
          />
        </Box>
      </Card>

      {/* Import CSV Modal */}
      <Dialog open={importDialogOpen} onClose={() => !importing && setImportDialogOpen(false)} maxWidth="xs" fullWidth>
        <DialogTitle sx={{ fontWeight: 600 }}>Import Data Karyawan (CSV)</DialogTitle>
        <DialogContent>
          <Box sx={{ mt: 1, display: 'flex', flexDirection: 'column', gap: 2 }}>
            <Button variant="outlined" color="info" onClick={handleDownloadTemplate} size="small" startIcon={<DownloadIcon />}>
              Unduh Template CSV
            </Button>
            <Typography variant="body2" color="text.secondary">
              Pilih file CSV dengan format sesuai template di atas.
            </Typography>
            <input type="file" accept=".csv" onChange={handleFileChange} disabled={importing} />
          </Box>
        </DialogContent>
        <DialogActions sx={{ px: 3, pb: 2 }}>
          <Button onClick={() => setImportDialogOpen(false)} disabled={importing}>
            Batal
          </Button>
          <Button variant="contained" onClick={handleImport} disabled={!importFile || importing} startIcon={importing ? <CircularProgress size={16} /> : null}>
            {importing ? 'Mengimpor...' : 'Upload & Import'}
          </Button>
        </DialogActions>
      </Dialog>

      {/* Toast Notification */}
      <Snackbar open={snackbar.open} autoHideDuration={4000} onClose={() => setSnackbar({ ...snackbar, open: false })}>
        <Alert severity={snackbar.severity} onClose={() => setSnackbar({ ...snackbar, open: false })}>
          {snackbar.message}
        </Alert>
      </Snackbar>
    </Box>
  );
};

export default EmployeeList;
