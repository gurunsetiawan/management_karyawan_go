import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import {
  TextField,
  Button,
  Paper,
  Typography,
  Box,
  CircularProgress,
  FormControl,
  InputLabel,
  MenuItem,
  Select,
  SelectChangeEvent,
  IconButton,
} from '@mui/material';
import { Save, ArrowBack } from '@mui/icons-material';
import { getEmployee, createEmployee, updateEmployee, Employee } from '../services/api';

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
    role: '',
    phone: '',
    alamat: '',
  });

  useEffect(() => {
    if (!isEdit) return;

    const fetchEmployee = async () => {
      try {
        const employee = await getEmployee(parseInt(id!));
        const { id: _, created_at, ...employeeData } = employee;
        setFormData(employeeData);
      } catch (err) {
        setError('Failed to fetch employee data');
        console.error(err);
      } finally {
        setLoading(false);
      }
    };

    fetchEmployee();
  }, [id, isEdit]);

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
      setError('All fields are required');
      return false;
    }

    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(email)) {
      setError('Please enter a valid email address');
      return false;
    }

    return true;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!validateForm()) return;

    try {
      setLoading(true);
      if (isEdit) {
        await updateEmployee(parseInt(id!), formData);
      } else {
        await createEmployee(formData);
      }
      navigate('/');
    } catch (err) {
      setError('Failed to save employee');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" mt={4}>
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Paper sx={{ p: 3, maxWidth: 800, margin: '0 auto' }}>
      <Box display="flex" alignItems="center" mb={3}>
        <IconButton onClick={() => navigate(-1)} sx={{ mr: 1 }}>
          <ArrowBack />
        </IconButton>
        <Typography variant="h5">
          {isEdit ? 'Edit Employee' : 'Add New Employee'}
        </Typography>
      </Box>

      {error && (
        <Box mb={2} p={1} bgcolor="error.light" borderRadius={1}>
          <Typography color="error">{error}</Typography>
        </Box>
      )}

      <form onSubmit={handleSubmit}>
        <Box mb={2}>
          <TextField
            fullWidth
            label="Name"
            name="name"
            value={formData.name}
            onChange={handleChange}
            margin="normal"
            required
          />
        </Box>

        <Box mb={2}>
          <TextField
            fullWidth
            label="Email"
            name="email"
            type="email"
            value={formData.email}
            onChange={handleChange}
            margin="normal"
            required
          />
        </Box>

        <Box mb={2}>
          <TextField
            fullWidth
            label="Position"
            name="position"
            value={formData.position}
            onChange={handleChange}
            margin="normal"
            required
          />
        </Box>

        <Box mb={2}>
          <FormControl fullWidth margin="normal" required>
            <InputLabel id="role-label">Role</InputLabel>
            <Select
              labelId="role-label"
              name="role"
              value={formData.role}
              label="Role"
              onChange={handleRoleChange}
            >
              <MenuItem value="">
                <em>Select Role</em>
              </MenuItem>
              <MenuItem value="Manager">Manager</MenuItem>
              <MenuItem value="Developer">Developer</MenuItem>
              <MenuItem value="Designer">Designer</MenuItem>
              <MenuItem value="HR">HR</MenuItem>
              <MenuItem value="Other">Other</MenuItem>
            </Select>
          </FormControl>
        </Box>

        <Box mb={2}>
          <TextField
            fullWidth
            label="Phone"
            name="phone"
            value={formData.phone}
            onChange={handleChange}
            margin="normal"
            required
          />
        </Box>

        <Box mb={3}>
          <TextField
            fullWidth
            label="Address"
            name="alamat"
            value={formData.alamat}
            onChange={handleChange}
            margin="normal"
            multiline
            rows={3}
            required
          />
        </Box>

        <Box display="flex" justifyContent="flex-end" gap={2}>
          <Button
            variant="outlined"
            onClick={() => navigate('/')}
            disabled={loading}
          >
            Cancel
          </Button>
          <Button
            type="submit"
            variant="contained"
            color="primary"
            startIcon={loading ? <CircularProgress size={20} /> : <Save />}
            disabled={loading}
          >
            {isEdit ? 'Update' : 'Save'}
          </Button>
        </Box>
      </form>
    </Paper>
  );
};

export default EmployeeForm;
