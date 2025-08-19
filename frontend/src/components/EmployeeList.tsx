import React, { useState, useEffect } from 'react';
import {
  Box, Button, IconButton, Paper, Table, TableBody, TableCell,
  TableContainer, TableHead, TableRow, Typography, TextField,
  TablePagination, Dialog, DialogTitle, DialogContent, DialogActions,
  AppBar, Toolbar, Tooltip, Snackbar, Alert, CircularProgress
} from '@mui/material';
import {
  Add as AddIcon, Edit as EditIcon, Delete as DeleteIcon,
  Person as PersonIcon, Work as WorkIcon, Email as EmailIcon,
  Phone as PhoneIcon, Search as SearchIcon
} from '@mui/icons-material';
import { DataGrid, GridColDef, GridSortModel, GridToolbar } from '@mui/x-data-grid';
import { getEmployees, deleteEmployee, Employee } from '../services/api';

const EmployeeList: React.FC = () => {
  const [employees, setEmployees] = useState<Employee[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [openModal, setOpenModal] = useState<boolean>(false);
  const [selectedEmployee, setSelectedEmployee] = useState<Employee | null>(null);
  const [snackbar, setSnackbar] = useState({ open: false, message: '', severity: 'success' });
  const [paginationModel, setPaginationModel] = useState({
    page: 0,
    pageSize: 10,
  });
  const [sortModel, setSortModel] = useState<GridSortModel>([
    { field: 'name', sort: 'asc' },
  ]);
  const [searchTerm, setSearchTerm] = useState('');

  useEffect(() => {
    fetchEmployees();
  }, []);

  const fetchEmployees = async () => {
    try {
      setLoading(true);
      const data = await getEmployees();
      setEmployees(data);
    } catch (error) {
      showSnackbar('Gagal memuat data karyawan', 'error');
      console.error('Error fetching employees:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (id: number) => {
    if (window.confirm('Apakah Anda yakin ingin menghapus karyawan ini?')) {
      try {
        await deleteEmployee(id);
        fetchEmployees();
        showSnackbar('Data karyawan berhasil dihapus', 'success');
      } catch (error) {
        showSnackbar('Gagal menghapus data karyawan', 'error');
        console.error('Error deleting employee:', error);
      }
    }
  };

  const handleEdit = (employee: Employee) => {
    setSelectedEmployee(employee);
    setOpenModal(true);
  };

  const handleCloseModal = () => {
    setOpenModal(false);
    setSelectedEmployee(null);
  };

  const showSnackbar = (message: string, severity: 'success' | 'error' | 'info' | 'warning') => {
    setSnackbar({ open: true, message, severity });
  };

  const handleCloseSnackbar = () => {
    setSnackbar({ ...snackbar, open: false });
  };

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
      width: 120,
      renderCell: (params) => (
        <Box>
          <Tooltip title="Edit">
            <IconButton
              color="primary"
              onClick={(e) => {
                e.stopPropagation();
                handleEdit(params.row);
              }}
            >
              <EditIcon />
            </IconButton>
          </Tooltip>
          <Tooltip title="Hapus">
            <IconButton
              color="error"
              onClick={(e) => {
                e.stopPropagation();
                handleDelete(params.row.id);
              }}
            >
              <DeleteIcon />
            </IconButton>
          </Tooltip>
        </Box>
      ),
    },
  ];

  // Filter employees based on search term
  const filteredEmployees = React.useMemo(() => {
    if (!searchTerm) return employees;
    
    const term = searchTerm.toLowerCase();
    return employees.filter(employee => 
      employee.name?.toLowerCase().includes(term) ||
      employee.position?.toLowerCase().includes(term) ||
      employee.email?.toLowerCase().includes(term) ||
      employee.phone?.toLowerCase().includes(term) ||
      employee.role?.toLowerCase().includes(term) ||
      employee.alamat?.toLowerCase().includes(term)
    );
  }, [employees, searchTerm]);

  return (
    <Box sx={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
      <AppBar position="static" color="default" elevation={1}>
        <Toolbar>
          <WorkIcon sx={{ mr: 2 }} />
          <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
            Manajemen Karyawan
          </Typography>
          <Button
            variant="contained"
            color="primary"
            startIcon={<AddIcon />}
            onClick={() => setOpenModal(true)}
          >
            Tambah Karyawan
          </Button>
        </Toolbar>
      </AppBar>

      <Box sx={{ p: 3, flex: 1, display: 'flex', flexDirection: 'column' }}>
        <Paper sx={{ p: 2, mb: 2, display: 'flex', alignItems: 'center' }}>
          <SearchIcon sx={{ color: 'action.active', mr: 1 }} />
          <TextField
            fullWidth
            variant="standard"
            placeholder="Cari karyawan..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            InputProps={{
              disableUnderline: true,
            }}
            sx={{
              '& .MuiInputBase-input': {
                py: 1,
              },
            }}
          />
        </Paper>

        <Paper sx={{ flex: 1, display: 'flex', flexDirection: 'column' }}>
          <Box sx={{ height: '100%', width: '100%' }}>
            <DataGrid
              rows={filteredEmployees}
              columns={columns}
              pageSizeOptions={[5, 10, 25]}
              paginationModel={paginationModel}
              onPaginationModelChange={setPaginationModel}
              sortModel={sortModel}
              onSortModelChange={setSortModel}
              loading={loading}
              disableRowSelectionOnClick
              slots={{
                toolbar: GridToolbar,
                loadingOverlay: LinearProgress as any,
              }}
              slotProps={{
                toolbar: {
                  showQuickFilter: true,
                },
              }}
              localeText={{
                // Rows
                noRowsLabel: 'Tidak ada data',
                noResultsOverlayLabel: 'Tidak ada hasil yang ditemukan.',
                
                // Toolbar
                toolbarDensity: 'Ukuran',
                toolbarDensityLabel: 'Ukuran',
                toolbarDensityCompact: 'Kecil',
                toolbarDensityStandard: 'Standar',
                toolbarDensityComfortable: 'Nyaman',
                toolbarFilters: 'Filter',
                toolbarFiltersLabel: 'Tampilkan filter',
                toolbarFiltersTooltipHide: 'Sembunyikan filter',
                toolbarFiltersTooltipShow: 'Tampilkan filter',
                toolbarExport: 'Ekspor',
                toolbarExportCSV: 'Unduh sebagai CSV',
                toolbarExportPrint: 'Cetak',
                
                // Filter
                toolbarQuickFilterPlaceholder: 'Cari...',
                
                // Pagination
                MuiTablePagination: {
                  labelRowsPerPage: 'Baris per halaman:',
                  labelDisplayedRows: ({ from, to, count }: { from: number; to: number; count: number }) => 
                    `${from}-${to} dari ${count !== -1 ? count : `lebih dari ${to}`}`,
                },
                
                // Column menu
                columnMenuLabel: 'Menu',
                columnMenuShowColumns: 'Tampilkan kolom',
                columnMenuFilter: 'Filter',
                columnMenuHideColumn: 'Sembunyikan',
                columnMenuUnsort: 'Batal urutkan',
                columnMenuSortAsc: 'Urutkan A ke Z',
                columnMenuSortDesc: 'Urutkan Z ke A',
                
                // Filter operators
                filterOperatorContains: 'mengandung',
                filterOperatorEquals: 'sama dengan',
                filterOperatorStartsWith: 'dimulai dengan',
                filterOperatorEndsWith: 'diakhiri dengan',
                filterOperatorIsEmpty: 'kosong',
                filterOperatorIsNotEmpty: 'tidak kosong',
                filterOperatorIsAnyOf: 'salah satu dari',
                
                // Column header titles
                columnHeaderSortIconLabel: 'Urutkan',
              } as any}
            />
          </Box>
        </Paper>
      </Box>

      {/* Add/Edit Modal */}
      <Dialog open={openModal} onClose={handleCloseModal} maxWidth="sm" fullWidth>
        <DialogTitle>
          {selectedEmployee ? 'Edit Karyawan' : 'Tambah Karyawan Baru'}
        </DialogTitle>
        <DialogContent>
          <Box component="form" sx={{ mt: 2 }}>
            <TextField
              margin="normal"
              required
              fullWidth
              label="Nama"
              name="name"
              autoFocus
              defaultValue={selectedEmployee?.name || ''}
              sx={{ mb: 2 }}
            />
            <TextField
              margin="normal"
              required
              fullWidth
              label="Jabatan"
              name="position"
              defaultValue={selectedEmployee?.position || ''}
              sx={{ mb: 2 }}
            />
            <TextField
              margin="normal"
              required
              fullWidth
              label="Email"
              name="email"
              type="email"
              defaultValue={selectedEmployee?.email || ''}
              sx={{ mb: 2 }}
            />
            <TextField
              margin="normal"
              required
              fullWidth
              label="Telepon"
              name="phone"
              defaultValue={selectedEmployee?.phone || ''}
            />
          </Box>
        </DialogContent>
        <DialogActions sx={{ p: 2 }}>
          <Button onClick={handleCloseModal} color="inherit">
            Batal
          </Button>
          <Button variant="contained" color="primary">
            {selectedEmployee ? 'Simpan Perubahan' : 'Tambah'}
          </Button>
        </DialogActions>
      </Dialog>

      <Snackbar
        open={snackbar.open}
        autoHideDuration={6000}
        onClose={handleCloseSnackbar}
        anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
      >
        <Alert
          onClose={handleCloseSnackbar}
          severity={snackbar.severity as 'success' | 'error' | 'info' | 'warning'}
          sx={{ width: '100%' }}
        >
          {snackbar.message}
        </Alert>
      </Snackbar>
    </Box>
  );
};

// Fallback component for loading overlay
const LinearProgress = () => (
  <div style={{ width: '100%', position: 'absolute', top: 0 }}>
    <div style={{ height: '4px', backgroundColor: '#1976d2' }} />
  </div>
);

export default EmployeeList;
