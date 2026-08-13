import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import {
  TextField,
  Button,
  Card,
  Typography,
  Box,
  CircularProgress,
  FormControl,
  InputLabel,
  MenuItem,
  Select,
  SelectChangeEvent,
  IconButton,
  Alert,
  Grid
} from '@mui/material';
import { Save, ArrowBack } from '@mui/icons-material';
import api, { Employee } from '../services/api';

const EmployeeForm: React.FC = () => {
  const { id } = useParams<{ id?: string }>();
  const isEdit = Boolean(id);
  const navigate = useNavigate();
  const [loading, setLoading] = useState(isEdit);
  const [error, setError] = useState<string | null>(null);
  const [formData, setFormData] = useState<Omit<Employee, 'id' | 'created_at' | 'updated_at'>>({
    name: '',
    email: '',
    position: '',
    role: 'karyawan',
    phone: '',
    alamat: '',
  });

  useEffect(() => {
    const fetchEmployee = async () => {
      if (!id) return;
      try {
        const employee = await api.getEmployee(parseInt(id));
        if (employee) {
          const { id: _, created_at, updated_at, ...employeeData } = employee;
          setFormData(employeeData);
        }
      } catch (err) {
        setError('Gagal mengambil data karyawan');
        console.error(err);
      } finally {
        setLoading(false);
      }
    };

    fetchEmployee();
  }, [id]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value,
    }));
  };

  const handleRoleChange = (e: SelectChangeEvent) => {
    setFormData(prev => ({
      ...prev,
      role: e.target.value,
    }));
  };

  const validateForm = (): boolean => {
    const { name, email, position, role, phone, alamat } = formData;
    if (!name.trim() || !email.trim() || !position.trim() || !role || !phone.trim() || !alamat.trim()) {
      setError('Semua kolom wajib diisi');
      return false;
    }

    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(email)) {
      setError('Format email tidak valid');
      return false;
    }

    return true;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    
    if (!validateForm()) return;

    try {
      setLoading(true);
      if (isEdit && id) {
        await api.updateEmployee(parseInt(id), formData);
      } else {
        const newEmployee = {
          ...formData,
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString()
        };
        await api.createEmployee(newEmployee);
      }
      navigate('/');
    } catch (err) {
      setError('Gagal menyimpan data karyawan');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" mt={6}>
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Box sx={{ maxWidth: 800, mx: 'auto' }}>
      <Box display="flex" alignItems="center" mb={3} gap={1}>
        <IconButton onClick={() => navigate(-1)} size="small" sx={{ border: '1px solid', borderColor: 'divider' }}>
          <ArrowBack fontSize="small" />
        </IconButton>
        <Typography variant="h5" sx={{ fontWeight: 700 }}>
          {isEdit ? 'Edit Karyawan' : 'Tambah Karyawan Baru'}
        </Typography>
      </Box>

      <Card sx={{ p: 4 }}>
        {error && (
          <Alert severity="error" sx={{ mb: 3 }} onClose={() => setError(null)}>
            {error}
          </Alert>
        )}

        <form onSubmit={handleSubmit}>
          <Grid container spacing={2.5}>
            <Grid size={{ xs: 12, sm: 6 }}>
              <TextField
                fullWidth
                label="Nama Lengkap"
                name="name"
                value={formData.name}
                onChange={handleChange}
                required
                size="small"
              />
            </Grid>

            <Grid size={{ xs: 12, sm: 6 }}>
              <TextField
                fullWidth
                label="Email"
                name="email"
                type="email"
                value={formData.email}
                onChange={handleChange}
                required
                size="small"
              />
            </Grid>

            <Grid size={{ xs: 12, sm: 6 }}>
              <TextField
                fullWidth
                label="Jabatan / Departemen"
                name="position"
                value={formData.position}
                onChange={handleChange}
                required
                size="small"
              />
            </Grid>

            <Grid size={{ xs: 12, sm: 6 }}>
              <FormControl fullWidth size="small" required>
                <InputLabel id="role-label">Role Akses</InputLabel>
                <Select
                  labelId="role-label"
                  name="role"
                  value={formData.role}
                  label="Role Akses"
                  onChange={handleRoleChange}
                >
                  <MenuItem value="karyawan">Karyawan</MenuItem>
                  <MenuItem value="manager">Manager</MenuItem>
                  <MenuItem value="admin">Admin</MenuItem>
                </Select>
              </FormControl>
            </Grid>

            <Grid size={{ xs: 12, sm: 6 }}>
              <TextField
                fullWidth
                label="No. Telepon"
                name="phone"
                value={formData.phone}
                onChange={handleChange}
                required
                size="small"
              />
            </Grid>

            <Grid size={{ xs: 12 }}>
              <TextField
                fullWidth
                label="Alamat"
                name="alamat"
                value={formData.alamat}
                onChange={handleChange}
                multiline
                rows={3}
                required
              />
            </Grid>
          </Grid>

          <Box display="flex" justifyContent="flex-end" gap={1.5} mt={4}>
            <Button
              variant="outlined"
              onClick={() => navigate('/')}
              disabled={loading}
              sx={{ px: 3 }}
            >
              Batal
            </Button>
            <Button
              type="submit"
              variant="contained"
              color="primary"
              startIcon={loading ? <CircularProgress size={16} /> : <Save />}
              disabled={loading}
              sx={{ px: 3 }}
            >
              {isEdit ? 'Simpan Perubahan' : 'Simpan Karyawan'}
            </Button>
          </Box>
        </form>
      </Card>
    </Box>
  );
};

export default EmployeeForm;
