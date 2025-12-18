import React, { useState } from 'react';
import { useNavigate, useLocation } from 'react-router';
import {
  Drawer,
  List,
  ListItem,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  Box,
  Typography,
  Divider,
  IconButton,
  Avatar,
  Menu,
  MenuItem,
} from '@mui/material';
import MenuIcon from '@mui/icons-material/Menu';
import PeopleIcon from '@mui/icons-material/People';
import CategoryIcon from '@mui/icons-material/Category';
import SchoolIcon from '@mui/icons-material/School';
import ClassIcon from '@mui/icons-material/Class';
import AssignmentIcon from '@mui/icons-material/Assignment';
import LogoutIcon from '@mui/icons-material/Logout';
import AccountCircleIcon from '@mui/icons-material/AccountCircle';
import ArrowBackIcon from '@mui/icons-material/ArrowBack';
import type { User } from '../../../domain/entities/User';
import { ROLE_LABELS } from '../../../shared/constants/roles';
import { ROUTES } from '../../../shared/constants/routes';

interface SidebarProps {
  user: User | null;
  onLogout: () => void;
}

const drawerWidth = 240;

export default function Sidebar({ user, onLogout }: SidebarProps) {
  const navigate = useNavigate();
  const location = useLocation();
  const [mobileOpen, setMobileOpen] = useState(false);
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);

  const handleDrawerToggle = () => {
    setMobileOpen(!mobileOpen);
  };

  const handleMenuOpen = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleMenuClose = () => {
    setAnchorEl(null);
  };

  const isAdmin = user?.role === 'admin';
  const isInstructor = user?.role === 'instructor' || isAdmin;
  const isStudent = user?.role === 'student';
  
  const isInClassView = location.pathname.startsWith('/dashboard/student/classes/');

  const menuItems = [
    ...(isAdmin
      ? [
          {
            text: 'Users',
            icon: <PeopleIcon />,
            path: ROUTES.ADMIN_USERS,
          },
          {
            text: 'Categories',
            icon: <CategoryIcon />,
            path: ROUTES.ADMIN_CATEGORIES,
          },
          {
            text: 'Courses',
            icon: <SchoolIcon />,
            path: ROUTES.ADMIN_COURSES,
          },
          {
            text: 'Course Offerings',
            icon: <ClassIcon />,
            path: ROUTES.ADMIN_COURSE_OFFERINGS,
          },
          {
            text: 'Enrollments',
            icon: <AssignmentIcon />,
            path: ROUTES.ADMIN_ENROLLMENTS,
          },
        ]
      : []),
    ...(isInstructor && !isAdmin
      ? [
          {
            text: 'Course Offerings',
            icon: <ClassIcon />,
            path: ROUTES.INSTRUCTOR_COURSE_OFFERINGS,
          },
          {
            text: 'Enrollments',
            icon: <AssignmentIcon />,
            path: ROUTES.INSTRUCTOR_ENROLLMENTS,
          },
        ]
      : []),
    ...(isStudent
      ? [
          isInClassView
            ? {
                text: 'Back to Classes',
                icon: <ArrowBackIcon />,
                path: ROUTES.STUDENT_DASHBOARD,
              }
            : {
                text: 'Classes',
                icon: <ClassIcon />,
                path: ROUTES.STUDENT_DASHBOARD,
          },
        ]
      : []),
  ];

  const drawer = (
    <Box>
      <Box
        sx={{
          p: 2,
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-between',
        }}
      >
        <Typography variant="h6" noWrap component="div" sx={{ fontWeight: 600 }}>
          ASTO LMS
        </Typography>
        <IconButton
          onClick={handleMenuOpen}
          sx={{ ml: 'auto' }}
          size="small"
        >
          <Avatar sx={{ width: 32, height: 32, bgcolor: 'primary.main' }}>
            {user?.username?.[0]?.toUpperCase() || 'U'}
          </Avatar>
        </IconButton>
      </Box>
      <Menu
        anchorEl={anchorEl}
        open={Boolean(anchorEl)}
        onClose={handleMenuClose}
        anchorOrigin={{
          vertical: 'bottom',
          horizontal: 'right',
        }}
        transformOrigin={{
          vertical: 'top',
          horizontal: 'right',
        }}
      >
        <MenuItem disabled>
          <Box>
            <Typography variant="body2" fontWeight={600}>
              {user?.username}
            </Typography>
            <Typography variant="caption" color="text.secondary">
              {user?.email}
            </Typography>
            <Typography variant="caption" display="block" color="text.secondary">
              {user?.role ? ROLE_LABELS[user.role] : ''}
            </Typography>
          </Box>
        </MenuItem>
        <Divider />
        <MenuItem onClick={onLogout}>
          <ListItemIcon>
            <LogoutIcon fontSize="small" />
          </ListItemIcon>
          <ListItemText>Logout</ListItemText>
        </MenuItem>
      </Menu>
      <Divider />
      <List>
        {menuItems.map((item) => (
          <ListItem key={item.path} disablePadding>
            <ListItemButton
              selected={location.pathname === item.path}
              onClick={() => {
                navigate(item.path);
                setMobileOpen(false);
              }}
            >
              <ListItemIcon
                sx={{
                  color: location.pathname === item.path ? 'primary.main' : 'inherit',
                }}
              >
                {item.icon}
              </ListItemIcon>
              <ListItemText primary={item.text} />
            </ListItemButton>
          </ListItem>
        ))}
      </List>
    </Box>
  );

  return (
    <Box
      component="nav"
      sx={{ width: { sm: drawerWidth }, flexShrink: { sm: 0 } }}
      aria-label="dashboard navigation"
    >
      <IconButton
        color="inherit"
        aria-label="open drawer"
        edge="start"
        onClick={handleDrawerToggle}
        sx={{ mr: 2, display: { sm: 'none' }, position: 'fixed', top: 16, left: 16, zIndex: 1300 }}
      >
        <MenuIcon />
      </IconButton>
      <Drawer
        variant="temporary"
        open={mobileOpen}
        onClose={handleDrawerToggle}
        ModalProps={{
          keepMounted: true,
        }}
        sx={{
          display: { xs: 'block', sm: 'none' },
          '& .MuiDrawer-paper': { boxSizing: 'border-box', width: drawerWidth },
        }}
      >
        {drawer}
      </Drawer>
      <Drawer
        variant="permanent"
        sx={{
          display: { xs: 'none', sm: 'block' },
          '& .MuiDrawer-paper': { boxSizing: 'border-box', width: drawerWidth },
        }}
        open
      >
        {drawer}
      </Drawer>
    </Box>
  );
}

