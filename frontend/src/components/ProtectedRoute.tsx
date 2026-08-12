import React from 'react';
import { Navigate, Outlet, useLocation } from 'react-router-dom';
import { getCurrentUser } from '../services/auth';

interface ProtectedRouteProps {
  requiredRole?: string;
  children?: React.ReactNode;
}

const ProtectedRoute: React.FC<ProtectedRouteProps> = ({ requiredRole, children }) => {
  const location = useLocation();
  const user = getCurrentUser();
  
  if (!user) {
    // Redirect to login page if not authenticated
    return <Navigate to="/login" state={{ from: location }} replace />;
  }

  // Check if user has required role (if specified)
  if (requiredRole && user.role !== requiredRole) {
    // Redirect to not-authorized page or home if user doesn't have required role
    return <Navigate to="/" state={{ from: location }} replace />;
  }

  // If there are children, render them, otherwise render the Outlet
  return children ? <>{children}</> : <Outlet />;
};

export default ProtectedRoute;
