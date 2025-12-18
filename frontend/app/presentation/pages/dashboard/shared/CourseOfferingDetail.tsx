import React, { useState, useCallback } from 'react';
import { useParams, useNavigate } from 'react-router';
import {
  Box,
  Typography,
  Card,
  CardContent,
  Button,
  IconButton,
  Accordion,
  AccordionSummary,
  AccordionDetails,
  Chip,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
} from '@mui/material';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import AddIcon from '@mui/icons-material/Add';
import EditIcon from '@mui/icons-material/Edit';
import DeleteIcon from '@mui/icons-material/Delete';
import ArrowBackIcon from '@mui/icons-material/ArrowBack';
import PeopleIcon from '@mui/icons-material/People';
import {
  useGetCourseOfferingQuery,
  useCreateCourseOfferingMutation,
  useUpdateCourseOfferingMutation,
  useDeleteCourseOfferingMutation,
  useAssignInstructorMutation,
  useRemoveInstructorMutation,
} from '../../../../infrastructure/api/rtk/courseOfferingsApi';
import { useGetUsersQuery } from '../../../../infrastructure/api/rtk/usersApi';
import { useSearchParams } from 'react-router';
import CourseOfferingForm from '../../../components/course/CourseOfferingForm';
import Success from '../../../components/common/Success';
import {
  useGetCourseSectionsQuery,
  useCreateCourseSectionMutation,
  useUpdateCourseSectionMutation,
  useDeleteCourseSectionMutation,
} from '../../../../infrastructure/api/rtk/courseSectionsApi';
import {
  useGetSectionModulesQuery,
  useCreateSectionModuleMutation,
  useUpdateSectionModuleMutation,
  useDeleteSectionModuleMutation,
} from '../../../../infrastructure/api/rtk/sectionModulesApi';
import { useDashboardAuth } from '../../../hooks/useDashboardAuth';
import DashboardLayout from '../../../components/dashboard/Layout';
import Loading from '../../../components/common/Loading';
import Error from '../../../components/common/Error';
import ModuleContentModal from '../../../components/courses/ModuleContentModal';
import type { SectionModule } from '../../../../domain/entities/SectionModule';

interface CourseOfferingDetailProps {
  routePrefix: string;
  allowedRoles: ('admin' | 'instructor')[];
}

export default function CourseOfferingDetail({ routePrefix, allowedRoles }: CourseOfferingDetailProps) {
  const { id } = useParams<{ id: string }>();
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const { user: currentUser, logout } = useDashboardAuth();
  const [expandedSection, setExpandedSection] = useState<string | false>(false);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);
  const [sectionDialogOpen, setSectionDialogOpen] = useState(false);
  const [moduleDialogOpen, setModuleDialogOpen] = useState(false);
  const [isEditing, setIsEditing] = useState(false);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [editingSection, setEditingSection] = useState<{ id: string; name: string; description: string; order: number; status: 'draft' | 'published' | 'archived' } | null>(null);
  const [editingModule, setEditingModule] = useState<{ id: string; sectionId: string; name: string; description: string; contentType: string; order: number } | null>(null);
  const [contentModalOpen, setContentModalOpen] = useState(false);
  const [contentModalModule, setContentModalModule] = useState<SectionModule | null>(null);
  const [contentModalMode, setContentModalMode] = useState<'create' | 'manage'>('create');
  const [instructorDialogOpen, setInstructorDialogOpen] = useState(false);
  const [selectedInstructorId, setSelectedInstructorId] = useState<string>('');
  const isNew = id === 'new';
  const isAdmin = currentUser?.role === 'admin';
  const courseIdFromQuery = searchParams.get('courseId');

  const { data: offering, isLoading: offeringLoading, error: offeringError } = useGetCourseOfferingQuery(id!, {
    skip: isNew,
  });
  const { data: sectionsData, isLoading: sectionsLoading } = useGetCourseSectionsQuery(
    { offeringId: id! },
    { skip: isNew }
  );
  const [createOffering, { isLoading: creating, error: createError }] = useCreateCourseOfferingMutation();
  const [updateOffering, { isLoading: updating, error: updateError }] = useUpdateCourseOfferingMutation();
  const [deleteOffering, { isLoading: deleting }] = useDeleteCourseOfferingMutation();
  const [createSection, { isLoading: creatingSection }] = useCreateCourseSectionMutation();
  const [updateSection, { isLoading: updatingSection }] = useUpdateCourseSectionMutation();
  const [deleteSection] = useDeleteCourseSectionMutation();
  const [createModule, { isLoading: creatingModule }] = useCreateSectionModuleMutation();
  const [updateModule, { isLoading: updatingModule }] = useUpdateSectionModuleMutation();
  const [deleteModule] = useDeleteSectionModuleMutation();
  const [assignInstructor, { isLoading: assigningInstructor }] = useAssignInstructorMutation();
  const [removeInstructor, { isLoading: removingInstructor }] = useRemoveInstructorMutation();
  const { data: instructorsData } = useGetUsersQuery({ role: 'instructor', limit: 100 });

  const sections = sectionsData?.sections || [];

  const handleSubmit = useCallback(
    async (data: {
      courseId: string;
      name: string;
      description: string;
      offeringType: 'online' | 'oncampus';
      duration?: string | null;
      classTime?: string | null;
      enrollmentCost: number;
      status?: 'pending' | 'active' | 'ongoing' | 'completed';
    }) => {
      try {
        if (isNew) {
          const result = await createOffering({
            courseId: data.courseId,
            data: {
              name: data.name,
              description: data.description,
              offeringType: data.offeringType,
              duration: data.duration,
              classTime: data.classTime,
              enrollmentCost: data.enrollmentCost,
            },
          }).unwrap();
          setSuccessMessage('Course offering created successfully');
          setTimeout(() => {
            navigate(`${routePrefix}/course-offerings/${result.id}`);
          }, 1500);
        } else {
          await updateOffering({
            id: id!,
            data: {
              name: data.name,
              description: data.description,
              offeringType: data.offeringType,
              duration: data.duration,
              classTime: data.classTime,
              enrollmentCost: data.enrollmentCost,
              status: data.status,
            },
          }).unwrap();
          setSuccessMessage('Course offering updated successfully');
          setIsEditing(false);
        }
      } catch (err) {
        console.error('Failed to save course offering:', err);
      }
    },
    [isNew, createOffering, updateOffering, id, navigate, routePrefix]
  );

  const handleSectionChange = (sectionId: string) => (_event: React.SyntheticEvent, isExpanded: boolean) => {
    setExpandedSection(isExpanded ? sectionId : false);
  };

  const handleAddSection = useCallback(() => {
    setEditingSection(null);
    setSectionDialogOpen(true);
  }, []);

  const handleEditSection = useCallback((section: { id: string; name: string; description: string; order: number; status: 'draft' | 'published' | 'archived' }) => {
    setEditingSection(section);
    setSectionDialogOpen(true);
  }, []);

  const handleAddModule = useCallback((sectionId: string) => {
    setEditingModule({ id: '', sectionId, name: '', description: '', contentType: 'zoom', order: 0 });
    setModuleDialogOpen(true);
  }, []);

  const handleEditModule = useCallback((module: { id: string; sectionId: string; name: string; description: string; contentType: string; order: number }) => {
    setEditingModule(module);
    setModuleDialogOpen(true);
  }, []);

  const handleDeleteSection = useCallback(async (sectionId: string) => {
    if (window.confirm('Are you sure you want to delete this section?')) {
      try {
        await deleteSection(sectionId).unwrap();
      } catch (err) {
        console.error('Failed to delete section:', err);
      }
    }
  }, [deleteSection]);

  const handleDeleteModule = useCallback(async (moduleId: string) => {
    if (window.confirm('Are you sure you want to delete this module?')) {
      try {
        await deleteModule(moduleId).unwrap();
      } catch (err) {
        console.error('Failed to delete module:', err);
      }
    }
  }, [deleteModule]);

  const handleDeleteOffering = useCallback(async () => {
    try {
      await deleteOffering(id!).unwrap();
      setSuccessMessage('Course offering deleted successfully');
      setTimeout(() => {
        navigate(`${routePrefix}/course-offerings`);
      }, 1500);
    } catch (err) {
      console.error('Failed to delete course offering:', err);
    }
  }, [deleteOffering, id, navigate, routePrefix]);

  const handleCreateContent = useCallback((module: SectionModule) => {
    setContentModalModule(module);
    setContentModalMode('create');
    setContentModalOpen(true);
  }, []);

  const handleManageContent = useCallback((module: SectionModule) => {
    setContentModalModule(module);
    setContentModalMode('manage');
    setContentModalOpen(true);
  }, []);

  const handleContentModalClose = useCallback(() => {
    setContentModalOpen(false);
    setContentModalModule(null);
  }, []);

  const handleAssignInstructor = useCallback(async () => {
    if (!selectedInstructorId || !offering) return;
    
    const instructor = instructorsData?.users.find(u => u.id === selectedInstructorId);
    if (!instructor || !instructor.id || !instructor.username) {
      console.error('Invalid instructor data:', instructor);
      return;
    }

    try {
      const requestData = {
        instructor_id: instructor.id,
        instructor_username: instructor.username,
      };
      
      await assignInstructor({
        offeringId: offering.id,
        data: requestData,
      }).unwrap();
      setSuccessMessage('Instructor assigned successfully');
      setInstructorDialogOpen(false);
      setSelectedInstructorId('');
    } catch (err) {
      console.error('Failed to assign instructor:', err);
    }
  }, [selectedInstructorId, offering, instructorsData, assignInstructor]);

  const handleRemoveInstructor = useCallback(async (instructorId: string) => {
    if (!offering || !window.confirm('Are you sure you want to remove this instructor?')) return;

    try {
      await removeInstructor({
        offeringId: offering.id,
        instructorId,
      }).unwrap();
      setSuccessMessage('Instructor removed successfully');
    } catch (err) {
      console.error('Failed to remove instructor:', err);
    }
  }, [offering, removeInstructor]);

  if (!currentUser || !allowedRoles.includes(currentUser.role as 'admin' | 'instructor')) {
    return (
      <Box sx={{ p: 4, textAlign: 'center' }}>
        <Typography variant="h5" color="error">
          Access Denied. {allowedRoles.map(r => r.charAt(0).toUpperCase() + r.slice(1)).join(' or ')} role required.
        </Typography>
      </Box>
    );
  }

  if (isNew) {
    return (
      <DashboardLayout user={currentUser} onLogout={logout}>
        <Box>
          <Box sx={{ display: 'flex', alignItems: 'center', mb: 3, gap: 2 }}>
            <IconButton onClick={() => navigate(`${routePrefix}/course-offerings`)}>
              <ArrowBackIcon />
            </IconButton>
            <Typography variant="h4" component="h1" sx={{ fontWeight: 600 }}>
              Create Course Offering
            </Typography>
          </Box>
          {successMessage && (
            <Success
              message={successMessage}
              autoDismiss
              autoDismissDelay={3000}
              onDismiss={() => setSuccessMessage(null)}
              fullWidth
            />
          )}
          <CourseOfferingForm
            courseId={courseIdFromQuery}
            onSubmit={handleSubmit}
            onCancel={() => navigate(`${routePrefix}/course-offerings`)}
            loading={creating}
            error={createError ? 'Failed to create course offering' : undefined}
          />
        </Box>
      </DashboardLayout>
    );
  }

  if (offeringLoading) {
    return (
      <DashboardLayout user={currentUser} onLogout={logout}>
        <Loading variant="skeleton" fullWidth />
      </DashboardLayout>
    );
  }

  if (offeringError || !offering) {
    return (
      <DashboardLayout user={currentUser} onLogout={logout}>
        <Error message="Failed to load course offering" onRetry={() => window.location.reload()} fullWidth />
      </DashboardLayout>
    );
  }

  return (
    <DashboardLayout user={currentUser} onLogout={logout}>
      <Box>
        <Box sx={{ display: 'flex', alignItems: 'center', mb: 3, gap: 2 }}>
          <IconButton onClick={() => navigate(`${routePrefix}/course-offerings`)}>
            <ArrowBackIcon />
          </IconButton>
          <Typography variant="h4" component="h1" sx={{ fontWeight: 600 }}>
            {offering.name}
          </Typography>
          {!isEditing && (
            <Box sx={{ ml: 'auto', display: 'flex', gap: 1 }}>
              <Button
                variant="outlined"
                startIcon={<EditIcon />}
                onClick={() => setIsEditing(true)}
              >
                Edit
              </Button>
              <Button
                variant="outlined"
                color="error"
                startIcon={<DeleteIcon />}
                onClick={() => setDeleteDialogOpen(true)}
              >
                Delete
              </Button>
            </Box>
          )}
        </Box>

        {successMessage && (
          <Success
            message={successMessage}
            autoDismiss
            autoDismissDelay={3000}
            onDismiss={() => setSuccessMessage(null)}
            fullWidth
          />
        )}

        {isEditing ? (
          <CourseOfferingForm
            courseOffering={offering}
            onSubmit={handleSubmit}
            onCancel={() => setIsEditing(false)}
            loading={updating}
            error={updateError ? 'Failed to update course offering' : undefined}
          />
        ) : (
          <>
            <Card sx={{ mb: 3 }}>
              <CardContent>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
                  <Typography variant="h6">Course Offering Details</Typography>
                  <Button
                    variant="contained"
                    startIcon={<PeopleIcon />}
                    onClick={() => navigate(`${routePrefix}/enrollments?courseOfferingId=${offering.id}`)}
                  >
                    Manage Enrollments
                  </Button>
                </Box>
                <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1 }}>
                  {offering.courseName && <Typography><strong>Course:</strong> {offering.courseName}</Typography>}
                  <Typography><strong>Description:</strong> {offering.description}</Typography>
                  <Typography><strong>Type:</strong> {offering.offeringType}</Typography>
                  <Typography><strong>Status:</strong> {offering.status}</Typography>
                  <Typography><strong>Enrollment Cost:</strong> ${offering.enrollmentCost.toFixed(2)}</Typography>
                  {offering.duration && <Typography><strong>Duration:</strong> {offering.duration}</Typography>}
                  {offering.classTime && <Typography><strong>Class Time:</strong> {offering.classTime}</Typography>}
                </Box>
              </CardContent>
            </Card>

            <Card sx={{ mb: 3 }}>
              <CardContent>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
                  <Typography variant="h6">Instructors</Typography>
                  {isAdmin && (
                    <Button
                      variant="contained"
                      startIcon={<AddIcon />}
                      onClick={() => setInstructorDialogOpen(true)}
                    >
                      Assign Instructor
                    </Button>
                  )}
                </Box>
                {offering.instructors && offering.instructors.length > 0 ? (
                  <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1 }}>
                    {offering.instructors.map((instructor) => (
                      <Box
                        key={instructor.id}
                        sx={{
                          display: 'flex',
                          justifyContent: 'space-between',
                          alignItems: 'center',
                          p: 2,
                          border: '1px solid #e0e0e0',
                          borderRadius: 1,
                        }}
                      >
                        <Box>
                          <Typography variant="body1" sx={{ fontWeight: 500 }}>
                            {instructor.instructorUsername}
                          </Typography>
                          <Typography variant="caption" color="text.secondary">
                            ID: {instructor.instructorId}
                          </Typography>
                        </Box>
                        {isAdmin && (
                          <IconButton
                            size="small"
                            color="error"
                            onClick={() => handleRemoveInstructor(instructor.instructorId)}
                            disabled={removingInstructor}
                          >
                            <DeleteIcon />
                          </IconButton>
                        )}
                      </Box>
                    ))}
                  </Box>
                ) : (
                  <Typography variant="body2" color="text.secondary" sx={{ textAlign: 'center', py: 2 }}>
                    No instructors assigned yet.
                  </Typography>
                )}
              </CardContent>
            </Card>
          </>
        )}

        {!isEditing && (
          <Card>
            <CardContent>
              <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
                <Typography variant="h6">Sections</Typography>
                <Button
                  variant="contained"
                  startIcon={<AddIcon />}
                  onClick={handleAddSection}
                >
                  Add Section
                </Button>
              </Box>

              {sectionsLoading && <Loading variant="skeleton" />}

              {!sectionsLoading && sections.length === 0 && (
                <Typography variant="body2" color="text.secondary" sx={{ textAlign: 'center', py: 4 }}>
                  No sections yet. Click "Add Section" to create one.
                </Typography>
              )}

              {!sectionsLoading && sections.map((section) => (
                <SectionAccordion
                  key={section.id}
                  section={section}
                  expanded={expandedSection === section.id}
                  onChange={handleSectionChange(section.id)}
                  onEdit={() => handleEditSection(section)}
                  onAddModule={() => handleAddModule(section.id)}
                  onDeleteSection={() => handleDeleteSection(section.id)}
                  onEditModule={handleEditModule}
                  onDeleteModule={handleDeleteModule}
                  onCreateContent={handleCreateContent}
                  onManageContent={handleManageContent}
                />
              ))}
            </CardContent>
          </Card>
        )}

        <ModuleContentModal
          open={contentModalOpen}
          onClose={handleContentModalClose}
          module={contentModalModule}
          mode={contentModalMode}
        />

        <SectionDialog
          open={sectionDialogOpen}
          onClose={() => {
            setSectionDialogOpen(false);
            setEditingSection(null);
          }}
          offeringId={id!}
          section={editingSection}
          onCreate={createSection}
          onUpdate={updateSection}
          loading={creatingSection || updatingSection}
          onSuccess={() => {
            setSectionDialogOpen(false);
            setEditingSection(null);
            setSuccessMessage(editingSection ? 'Section updated successfully' : 'Section created successfully');
          }}
        />

        <ModuleDialog
          open={moduleDialogOpen}
          onClose={() => {
            setModuleDialogOpen(false);
            setEditingModule(null);
          }}
          module={editingModule}
          onCreate={createModule}
          onUpdate={updateModule}
          loading={creatingModule || updatingModule}
          onSuccess={() => {
            setModuleDialogOpen(false);
            setEditingModule(null);
            setSuccessMessage(editingModule?.id ? 'Module updated successfully' : 'Module created successfully');
          }}
        />

        <Dialog
          open={deleteDialogOpen}
          onClose={() => setDeleteDialogOpen(false)}
        >
          <DialogTitle>Delete Course Offering</DialogTitle>
          <DialogContent>
            <Typography>
              Are you sure you want to delete "{offering.name}"? This action cannot be undone.
            </Typography>
          </DialogContent>
          <DialogActions>
            <Button onClick={() => setDeleteDialogOpen(false)} disabled={deleting}>
              Cancel
            </Button>
            <Button
              onClick={handleDeleteOffering}
              color="error"
              variant="contained"
              disabled={deleting}
            >
              {deleting ? 'Deleting...' : 'Delete'}
            </Button>
          </DialogActions>
        </Dialog>

        <Dialog
          open={instructorDialogOpen}
          onClose={() => {
            setInstructorDialogOpen(false);
            setSelectedInstructorId('');
          }}
          maxWidth="sm"
          fullWidth
        >
          <DialogTitle>Assign Instructor</DialogTitle>
          <DialogContent>
            <FormControl fullWidth margin="dense" sx={{ mt: 1 }}>
              <InputLabel>Select Instructor</InputLabel>
              <Select
                value={selectedInstructorId}
                label="Select Instructor"
                onChange={(e) => setSelectedInstructorId(e.target.value)}
              >
                {instructorsData?.users
                  .filter(instructor => 
                    !offering?.instructors?.some(
                      assigned => assigned.instructorId === instructor.id
                    )
                  )
                  .map((instructor) => (
                    <MenuItem key={instructor.id} value={instructor.id}>
                      {instructor.username} ({instructor.email})
                    </MenuItem>
                  ))}
              </Select>
            </FormControl>
          </DialogContent>
          <DialogActions>
            <Button
              onClick={() => {
                setInstructorDialogOpen(false);
                setSelectedInstructorId('');
              }}
              disabled={assigningInstructor}
            >
              Cancel
            </Button>
            <Button
              onClick={handleAssignInstructor}
              variant="contained"
              disabled={assigningInstructor || !selectedInstructorId}
            >
              {assigningInstructor ? 'Assigning...' : 'Assign'}
            </Button>
          </DialogActions>
        </Dialog>
      </Box>
    </DashboardLayout>
  );
}

function SectionAccordion({
  section,
  expanded,
  onChange,
  onEdit,
  onAddModule,
  onDeleteSection,
  onEditModule,
  onDeleteModule,
  onCreateContent,
  onManageContent,
}: {
  section: { id: string; name: string; description: string; order: number; status: string };
  expanded: boolean;
  onChange: (event: React.SyntheticEvent, isExpanded: boolean) => void;
  onEdit: () => void;
  onAddModule: () => void;
  onDeleteSection: () => void;
  onEditModule: (module: { id: string; sectionId: string; name: string; description: string; contentType: string; order: number }) => void;
  onDeleteModule: (moduleId: string) => void;
  onCreateContent: (module: SectionModule) => void;
  onManageContent: (module: SectionModule) => void;
}) {
  const { data: modulesData, isLoading: modulesLoading } = useGetSectionModulesQuery({ sectionId: section.id });
  const modules = modulesData?.modules || [];

  return (
    <Accordion expanded={expanded} onChange={onChange}>
      <AccordionSummary expandIcon={<ExpandMoreIcon />}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', width: '100%', mr: 2 }}>
          <Box>
            <Typography variant="subtitle1" sx={{ fontWeight: 500 }}>
              {section.name} (Order: {section.order})
            </Typography>
            <Chip label={section.status} size="small" sx={{ mt: 0.5 }} />
          </Box>
          <Box sx={{ display: 'flex', gap: 0.5 }}>
            <IconButton
              size="small"
              onClick={(e) => {
                e.stopPropagation();
                onEdit();
              }}
              color="primary"
            >
              <EditIcon />
            </IconButton>
            <IconButton
              size="small"
              onClick={(e) => {
                e.stopPropagation();
                onDeleteSection();
              }}
              color="error"
            >
              <DeleteIcon />
            </IconButton>
          </Box>
        </Box>
      </AccordionSummary>
      <AccordionDetails>
        <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
          {section.description}
        </Typography>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
          <Typography variant="subtitle2">Modules ({modules.length})</Typography>
          <Button
            size="small"
            variant="outlined"
            startIcon={<AddIcon />}
            onClick={onAddModule}
          >
            Add Module
          </Button>
        </Box>
        {modulesLoading && <Loading variant="skeleton" />}
        {!modulesLoading && modules.length === 0 && (
          <Typography variant="body2" color="text.secondary" sx={{ textAlign: 'center', py: 2 }}>
            No modules yet.
          </Typography>
        )}
        {!modulesLoading && modules.map((module) => {
          const hasContent = module.contentId && module.contentStatus === 'created';
          const needsContent = module.contentType === 'zoom' && !hasContent;
          
          return (
            <Box
              key={module.id}
              sx={{
                display: 'flex',
                justifyContent: 'space-between',
                alignItems: 'center',
                p: 2,
                mb: 1,
                border: '1px solid #e0e0e0',
                borderRadius: 1,
              }}
            >
              <Box>
                <Typography variant="body2" sx={{ fontWeight: 500 }}>
                  {module.name} (Order: {module.order})
                </Typography>
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mt: 0.5 }}>
                  <Typography variant="caption" color="text.secondary">
                    Type: {module.contentType}
                  </Typography>
                  <Chip
                    label={module.contentStatus}
                    size="small"
                    color={
                      module.contentStatus === 'created'
                        ? 'success'
                        : module.contentStatus === 'pending'
                        ? 'warning'
                        : 'default'
                    }
                  />
                </Box>
              </Box>
              <Box sx={{ display: 'flex', gap: 0.5 }}>
                {module.contentType === 'zoom' && (
                  needsContent ? (
                    <Button
                      size="small"
                      variant="outlined"
                      color="primary"
                      onClick={() => onCreateContent(module)}
                      sx={{ mr: 0.5 }}
                    >
                      Create Meeting
                    </Button>
                  ) : (
                    <Button
                      size="small"
                      variant="outlined"
                      color="primary"
                      onClick={() => onManageContent(module)}
                      sx={{ mr: 0.5 }}
                    >
                      Manage Content
                    </Button>
                  )
                )}
                <IconButton
                  size="small"
                  onClick={() => onEditModule({
                    id: module.id,
                    sectionId: module.courseSectionId,
                    name: module.name,
                    description: module.description,
                    contentType: module.contentType,
                    order: module.order,
                  })}
                  color="primary"
                >
                  <EditIcon />
                </IconButton>
                <IconButton
                  size="small"
                  onClick={() => onDeleteModule(module.id)}
                  color="error"
                >
                  <DeleteIcon />
                </IconButton>
              </Box>
            </Box>
          );
        })}
      </AccordionDetails>
    </Accordion>
  );
}

function SectionDialog({
  open,
  onClose,
  offeringId,
  section,
  onCreate,
  onUpdate,
  loading,
  onSuccess,
}: {
  open: boolean;
  onClose: () => void;
  offeringId: string;
  section: { id: string; name: string; description: string; order: number; status: 'draft' | 'published' | 'archived' } | null;
  onCreate: (params: { offeringId: string; data: { name: string; description?: string; order?: number; status?: 'draft' | 'published' | 'archived' } }) => Promise<any>;
  onUpdate: (params: { id: string; data: { name?: string; description?: string; order?: number; status?: 'draft' | 'published' | 'archived' } }) => Promise<any>;
  loading: boolean;
  onSuccess: () => void;
}) {
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [order, setOrder] = useState<number>(0);
  const [status, setStatus] = useState<'draft' | 'published' | 'archived'>('draft');

  React.useEffect(() => {
    if (section) {
      setName(section.name);
      setDescription(section.description || '');
      setOrder(section.order);
      setStatus(section.status);
    } else {
      setName('');
      setDescription('');
      setOrder(0);
      setStatus('draft');
    }
  }, [section, open]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      if (section) {
        await onUpdate({
          id: section.id,
          data: {
            name,
            description: description || undefined,
            order,
            status,
          },
        });
      } else {
        await onCreate({
          offeringId,
          data: {
            name,
            description: description || undefined,
            order,
            status,
          },
        });
      }
      onSuccess();
    } catch (err) {
      console.error('Failed to save section:', err);
    }
  };

  return (
    <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth>
      <form onSubmit={handleSubmit}>
        <DialogTitle>{section ? 'Edit Section' : 'Create Section'}</DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            label="Name"
            fullWidth
            variant="outlined"
            value={name}
            onChange={(e) => setName(e.target.value)}
            required
            sx={{ mb: 2 }}
          />
          <TextField
            margin="dense"
            label="Description"
            fullWidth
            multiline
            rows={4}
            variant="outlined"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            sx={{ mb: 2 }}
          />
          <TextField
            margin="dense"
            label="Order"
            type="number"
            fullWidth
            variant="outlined"
            value={order}
            onChange={(e) => setOrder(parseInt(e.target.value) || 0)}
            inputProps={{ min: 0 }}
            sx={{ mb: 2 }}
          />
          <FormControl fullWidth margin="dense">
            <InputLabel>Status</InputLabel>
            <Select
              value={status}
              label="Status"
              onChange={(e) => setStatus(e.target.value as 'draft' | 'published' | 'archived')}
            >
              <MenuItem value="draft">Draft</MenuItem>
              <MenuItem value="published">Published</MenuItem>
              <MenuItem value="archived">Archived</MenuItem>
            </Select>
          </FormControl>
        </DialogContent>
        <DialogActions>
          <Button onClick={onClose} disabled={loading}>
            Cancel
          </Button>
          <Button type="submit" variant="contained" disabled={loading || !name.trim()}>
            {loading ? 'Saving...' : section ? 'Update' : 'Create'}
          </Button>
        </DialogActions>
      </form>
    </Dialog>
  );
}

function ModuleDialog({
  open,
  onClose,
  module,
  onCreate,
  onUpdate,
  loading,
  onSuccess,
}: {
  open: boolean;
  onClose: () => void;
  module: { id: string; sectionId: string; name: string; description: string; contentType: string; order: number } | null;
  onCreate: (params: { sectionId: string; data: { name: string; description?: string; contentType: 'zoom'; order?: number } }) => Promise<any>;
  onUpdate: (params: { id: string; data: { name?: string; description?: string; order?: number } }) => Promise<any>;
  loading: boolean;
  onSuccess: () => void;
}) {
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [contentType, setContentType] = useState<'zoom'>('zoom');

  React.useEffect(() => {
    if (module) {
      setName(module.name);
      setDescription(module.description || '');
      setContentType(module.contentType as 'zoom');
    } else {
      setName('');
      setDescription('');
      setContentType('zoom');
    }
  }, [module, open]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      if (module && module.id && module.id.trim() !== '') {
        await onUpdate({
          id: module.id,
          data: {
            name,
            description: description || undefined,
          },
        });
      } else if (module && module.sectionId) {
        await onCreate({
          sectionId: module.sectionId,
          data: {
            name,
            description: description || undefined,
            contentType: 'zoom',
          },
        });
      }
      onSuccess();
    } catch (err) {
      console.error('Failed to save module:', err);
    }
  };

  return (
    <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth>
      <form onSubmit={handleSubmit}>
        <DialogTitle>{module && module.id && module.id.trim() !== '' ? 'Edit Module' : 'Create Module'}</DialogTitle>
        <DialogContent>
          <FormControl fullWidth margin="dense" sx={{ mb: 2 }}>
            <InputLabel>Content Type</InputLabel>
            <Select
              value={contentType}
              label="Content Type"
              onChange={(e) => setContentType(e.target.value as 'zoom')}
              disabled={!!module?.id}
            >
              <MenuItem value="zoom">Zoom</MenuItem>
            </Select>
          </FormControl>
          <TextField
            autoFocus
            margin="dense"
            label="Name"
            fullWidth
            variant="outlined"
            value={name}
            onChange={(e) => setName(e.target.value)}
            required
            sx={{ mb: 2 }}
          />
          <TextField
            margin="dense"
            label="Description"
            fullWidth
            multiline
            rows={4}
            variant="outlined"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={onClose} disabled={loading}>
            Cancel
          </Button>
          <Button type="submit" variant="contained" disabled={loading || !name.trim()}>
            {loading ? 'Saving...' : (module && module.id && module.id.trim() !== '' ? 'Update' : 'Create')}
          </Button>
        </DialogActions>
      </form>
    </Dialog>
  );
}

