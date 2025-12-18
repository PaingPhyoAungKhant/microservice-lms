import React, { useState } from 'react';
import { useParams, useNavigate } from 'react-router';
import {
  Box,
  Typography,
  Grid,
  Card,
  CardContent,
  Container,
  IconButton,
  Chip,
  Accordion,
  AccordionSummary,
  AccordionDetails,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button as MUIButton,
} from '@mui/material';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import ArrowBackIcon from '@mui/icons-material/ArrowBack';
import { useCourseDetail } from '../hooks/useCourses';
import { useGetCourseOfferingsQuery } from '../../infrastructure/api/rtk/courseOfferingsApi';
import { useGetCourseSectionsQuery } from '../../infrastructure/api/rtk/courseSectionsApi';
import { useGetSectionModulesQuery } from '../../infrastructure/api/rtk/sectionModulesApi';
import { useCreateEnrollmentMutation } from '../../infrastructure/api/rtk/enrollmentsApi';
import { useAuth } from '../hooks/useAuth';
import { getFileDownloadUrl } from '../../infrastructure/api/utils';
import Loading from '../components/common/Loading';
import Error from '../components/common/Error';
import Success from '../components/common/Success';
import Button from '../components/common/Button';
import { ROUTES } from '../../shared/constants/routes';
import type { CourseOffering } from '../../domain/entities/CourseOffering';
import type { CourseSection } from '../../domain/entities/CourseSection';
import type { SectionModule } from '../../domain/entities/SectionModule';

export default function CourseDetail() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { user } = useAuth();
  const [expandedOfferingId, setExpandedOfferingId] = useState<string | false>(false);
  const [authDialogOpen, setAuthDialogOpen] = useState(false);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);
  const [enrollingOfferingId, setEnrollingOfferingId] = useState<string | null>(null);
  const [successDialogOpen, setSuccessDialogOpen] = useState(false);
  const { course, loading, error, refetch } = useCourseDetail(id);

  const {
    data: offeringsData,
    isLoading: offeringsLoading,
    error: offeringsError,
    refetch: refetchOfferings,
  } = useGetCourseOfferingsQuery(
    id
      ? {
          courseId: id,
          limit: 20,
          offset: 0,
          sortColumn: 'created_at',
          sortDirection: 'asc',
        }
      : undefined,
    {
      skip: !id,
    }
  );

  const [createEnrollment, { isLoading: creatingEnrollment }] = useCreateEnrollmentMutation();

  const handleBack = () => {
    navigate(ROUTES.COURSES);
  };

  const handleEnroll = async (offering: CourseOffering) => {
    if (!course) return;

    if (!user || user.role !== 'student') {
      setAuthDialogOpen(true);
      return;
    }

    try {
      setEnrollingOfferingId(offering.id);
      await createEnrollment({
        studentId: user.id,
        studentUsername: user.username,
        courseId: course.id,
        courseName: course.name,
        courseOfferingId: offering.id,
        courseOfferingName: offering.name,
      }).unwrap();
      setSuccessMessage('Enrollment request submitted successfully');
      setSuccessDialogOpen(true);
    } catch (e: any) {
      const message =
        (e && typeof e === 'object' && 'data' in e && (e as any).data?.message) ||
        'Failed to enroll. You might already be enrolled.';
      setSuccessMessage(message);
    } finally {
      setEnrollingOfferingId(null);
    }
  };

  const handleAuthDialogClose = () => {
    setAuthDialogOpen(false);
  };

  const handleGoToLogin = () => {
    navigate(ROUTES.LOGIN);
  };

  const handleGoToRegister = () => {
    navigate(ROUTES.REGISTER);
  };

  const handleSuccessDialogClose = () => {
    setSuccessDialogOpen(false);
    navigate(ROUTES.STUDENT_DASHBOARD);
  };

  if (loading) {
    return <Loading variant="fullscreen" />;
  }

  if (error) {
    return (
      <Container>
        <Error
          message={error.message || 'Failed to load course details'}
          onRetry={() => refetch()}
          fullWidth
        />
      </Container>
    );
  }

  if (!course) {
    return (
      <Container>
        <Box sx={{ textAlign: 'center', py: 8 }}>
          <Typography variant="h5" color="text.secondary">
            Course not found
          </Typography>
        </Box>
      </Container>
    );
  }

  const thumbnailUrl =
    course.thumbnailUrl ||
    (course.thumbnailId ? getFileDownloadUrl('course-thumbnails', course.thumbnailId) : '/placeholder-course.jpg');

  const activeOfferings: CourseOffering[] =
    offeringsData?.courseOfferings.filter((offering) => offering.status === 'active') || [];

  return (
    <Container maxWidth="lg" sx={{ py: 4, px: { xs: 2, md: 4 } }}>
      <Box sx={{ mb: 4, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <IconButton
          onClick={handleBack}
          sx={{
            color: 'text.secondary',
            fontSize: '1.25rem',
            fontWeight: 600,
            '&:hover': {
              color: 'text.primary',
              transform: 'scale(1.05)',
            },
            transition: 'all 0.3s',
          }}
        >
          <ArrowBackIcon sx={{ fontSize: 40, mr: 1 }} />
          <Typography variant="h6">Back to Courses</Typography>
        </IconButton>
        {course.categories && course.categories.length > 0 && (
          <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 1 }}>
            {course.categories.map((category) => (
              <Chip key={category.id} label={category.name} color="primary" variant="outlined" size="small" />
            ))}
          </Box>
        )}
      </Box>

      <Box
        sx={{
          borderRadius: 2,
          overflow: 'hidden',
          boxShadow: 3,
          mb: 4,
        }}
      >
            <Box
              sx={{
            pt: '40%',
                width: '100%',
            backgroundImage: `url('${thumbnailUrl}')`,
                backgroundSize: 'cover',
                backgroundPosition: 'center',
                backgroundRepeat: 'no-repeat',
              }}
            />
      </Box>

      <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
        <Typography
          variant="h3"
          sx={{ color: 'primary.main', fontSize: { xs: '2rem', md: '2.5rem' }, fontWeight: 700, mb: 1 }}
        >
          {course.name}
            </Typography>

        <Typography variant="h4" sx={{ mt: 4, color: 'primary.main', fontSize: '1.875rem' }}>
          Available Courses
        </Typography>
            {offeringsError && (
              <Error
                message="Failed to load course offerings"
                onRetry={() => refetchOfferings()}
                fullWidth
              />
            )}
            {offeringsLoading && <Loading variant="skeleton" fullWidth />}
            {!offeringsLoading && !offeringsError && activeOfferings.length === 0 && (
              <Box sx={{ py: 4 }}>
                <Typography variant="body2" color="text.secondary">
                  No active offerings available for this course yet.
            </Typography>
              </Box>
            )}
        {!offeringsLoading &&
          !offeringsError &&
          activeOfferings.map((offering) => (
            <Accordion
              key={offering.id}
              expanded={expandedOfferingId === offering.id}
              onChange={(_, isExpanded) =>
                setExpandedOfferingId(isExpanded ? offering.id : false)
              }
            >
              <AccordionSummary expandIcon={<ExpandMoreIcon />}>
                <Box
                  sx={{
                    display: 'flex',
                    flexDirection: 'column',
                    gap: 1,
                    width: '100%',
                  }}
                >
                  <Box
                    sx={{
                      display: 'flex',
                      justifyContent: 'space-between',
                      alignItems: 'center',
                      gap: 2,
                      flexWrap: 'wrap',
                    }}
                  >
                    <Typography variant="h6">{offering.name}</Typography>
                    <Box sx={{ display: 'flex', gap: 1, alignItems: 'center', flexWrap: 'wrap' }}>
                      <Chip label={offering.offeringType} size="small" />
                      <Chip label="Active" color="success" size="small" />
                      <Typography variant="body2" sx={{ fontWeight: 600 }}>
                        ${offering.enrollmentCost.toFixed(2)}
                      </Typography>
                      <Button
                        variant="fill"
                        size="sm"
                        onClick={() => handleEnroll(offering)}
                        disabled={creatingEnrollment && enrollingOfferingId === offering.id}
                      >
                        {creatingEnrollment && enrollingOfferingId === offering.id
                          ? 'Enrolling...'
                          : 'Enroll'}
                      </Button>
                    </Box>
                  </Box>
                  {offering.description && (
                    <Typography variant="body2" color="text.secondary">
                      {offering.description}
                    </Typography>
                  )}
                </Box>
              </AccordionSummary>
              <AccordionDetails>
                <OfferingSections offeringId={offering.id} />
              </AccordionDetails>
            </Accordion>
          ))}

        <Card sx={{ mt: 4 }}>
          <CardContent sx={{ display: 'flex', flexDirection: 'column', gap: 1 }}>
            <Typography variant="h5" sx={{ fontSize: '1.5rem', fontWeight: 600 }}>
              About this course
            </Typography>
            <Typography variant="body2" color="text.secondary">
              {course.description}
                      </Typography>
                    </CardContent>
                  </Card>
      </Box>

      <Dialog open={authDialogOpen} onClose={handleAuthDialogClose}>
        <DialogTitle>Login or Register to Enroll</DialogTitle>
        <DialogContent>
          <Typography variant="body2" sx={{ mb: 2 }}>
            You need to be logged in as a student to enroll in this course. Please login or create a
            student account first.
          </Typography>
        </DialogContent>
        <DialogActions>
          <MUIButton onClick={handleAuthDialogClose}>Close</MUIButton>
          <MUIButton onClick={handleGoToRegister} color="primary">
            Register
          </MUIButton>
          <MUIButton onClick={handleGoToLogin} variant="contained">
            Login
          </MUIButton>
        </DialogActions>
      </Dialog>

      <Dialog open={successDialogOpen} onClose={handleSuccessDialogClose}>
        <DialogTitle>Enrollment successful</DialogTitle>
        <DialogContent>
          <Typography variant="body2" sx={{ mb: 2 }}>
            Your enrollment has been submitted successfully. You can now access your classes from
            your student dashboard.
          </Typography>
        </DialogContent>
        <DialogActions>
          <MUIButton onClick={handleSuccessDialogClose} variant="contained">
            Go to My Classes
          </MUIButton>
        </DialogActions>
      </Dialog>
    </Container>
  );
}

function OfferingSections({ offeringId }: { offeringId: string }) {
  const { data: sectionsData, isLoading: sectionsLoading, error: sectionsError } =
    useGetCourseSectionsQuery({ offeringId });

  const sections: CourseSection[] =
    sectionsData?.sections.slice().sort((a, b) => a.order - b.order) || [];

  if (sectionsLoading) {
    return <Loading variant="skeleton" fullWidth />;
  }

  if (sectionsError) {
    return (
      <Error
        message="Failed to load sections"
        onRetry={undefined}
        fullWidth
      />
    );
  }

  if (sections.length === 0) {
    return (
      <Typography variant="body2" color="text.secondary">
        No sections have been published for this offering yet.
      </Typography>
    );
  }

  return (
    <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
      {sections.map((section) => (
        <Card key={section.id} variant="outlined">
          <CardContent>
            <Typography variant="subtitle1" sx={{ fontWeight: 600 }}>
              {section.name}
            </Typography>
            {section.description && (
              <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>
                {section.description}
              </Typography>
            )}
            <SectionModules sectionId={section.id} />
          </CardContent>
        </Card>
      ))}
            </Box>
  );
}

function SectionModules({ sectionId }: { sectionId: string }) {
  const { data: modulesData, isLoading: modulesLoading, error: modulesError } =
    useGetSectionModulesQuery({ sectionId });

  const modules: SectionModule[] =
    modulesData?.modules.slice().sort((a, b) => a.order - b.order) || [];

  if (modulesLoading) {
    return <Loading variant="skeleton" fullWidth />;
  }

  if (modulesError) {
    return (
      <Error
        message="Failed to load modules"
        onRetry={undefined}
        fullWidth
      />
    );
  }

  if (modules.length === 0) {
    return (
      <Typography variant="body2" color="text.secondary">
        Modules will appear here once published.
      </Typography>
    );
  }

  return (
    <Box sx={{ mt: 1, display: 'flex', flexDirection: 'column', gap: 1 }}>
      {modules.map((module) => (
        <Box
          key={module.id}
          sx={{
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'center',
            p: 1.5,
            borderRadius: 1,
            bgcolor: 'grey.50',
          }}
        >
          <Box>
            <Typography variant="body2" sx={{ fontWeight: 500 }}>
              {module.name}
            </Typography>
            {module.description && (
              <Typography variant="caption" color="text.secondary">
                {module.description}
              </Typography>
            )}
            </Box>
          <Box sx={{ display: 'flex', gap: 1, alignItems: 'center' }}>
            <Chip label={module.contentType} size="small" />
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
      ))}
    </Box>
  );
}


