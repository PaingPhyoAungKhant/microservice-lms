import React, { useMemo, useState, useEffect, useCallback } from 'react';
import { useParams, useNavigate } from 'react-router';
import { useDispatch } from 'react-redux';
import {
  Box,
  Typography,
  IconButton,
  Drawer,
  List,
  ListItemButton,
  ListItemText,
  Divider,
  useMediaQuery,
  Chip,
  Button as MUIButton,
  Card,
  CardContent,
  CircularProgress,
  Alert,
} from '@mui/material';
import MenuIcon from '@mui/icons-material/Menu';
import ArrowBackIcon from '@mui/icons-material/ArrowBack';
import VideoCallIcon from '@mui/icons-material/VideoCall';
import DownloadIcon from '@mui/icons-material/Download';
import PlayCircleOutlineIcon from '@mui/icons-material/PlayCircleOutline';
import CheckCircleOutlineIcon from '@mui/icons-material/CheckCircleOutline';
import { useTheme } from '@mui/material/styles';
import { useAuth } from '../../../hooks/useAuth';
import { clearPublicAuth, clearDashboardAuth } from '../../../../infrastructure/store/authSlice';
import { useGetCourseOfferingQuery } from '../../../../infrastructure/api/rtk/courseOfferingsApi';
import { useGetCourseSectionsQuery } from '../../../../infrastructure/api/rtk/courseSectionsApi';
import { useGetSectionModulesQuery, useLazyGetSectionModulesQuery } from '../../../../infrastructure/api/rtk/sectionModulesApi';
import { useGetZoomMeetingByModuleQuery } from '../../../../infrastructure/api/rtk/zoomMeetingsApi';
import { useListZoomRecordingsQuery } from '../../../../infrastructure/api/rtk/zoomRecordingsApi';
import Loading from '../../../components/common/Loading';
import Error from '../../../components/common/Error';
import { downloadFileWithAuth, getErrorMessage } from '../../../../infrastructure/api/utils';
import { ROUTES } from '../../../../shared/constants/routes';
import type { CourseSection } from '../../../../domain/entities/CourseSection';
import type { SectionModule } from '../../../../domain/entities/SectionModule';

export default function ClassView() {
  const { offeringId } = useParams<{ offeringId: string }>();
  const navigate = useNavigate();
  const dispatch = useDispatch();
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('md'));

  const { user, loading: authLoading } = useAuth();
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const [selectedSectionId, setSelectedSectionId] = useState<string | null>(null);
  const [selectedModuleId, setSelectedModuleId] = useState<string | null>(null);
  const [downloadingFileId, setDownloadingFileId] = useState<string | null>(null);
  const [downloadError, setDownloadError] = useState<string | null>(null);

  const {
    data: offering,
    isLoading: offeringLoading,
    error: offeringError,
  } = useGetCourseOfferingQuery(offeringId || '', {
    skip: !offeringId,
  });

  const {
    data: sectionsData,
    isLoading: sectionsLoading,
    error: sectionsError,
  } = useGetCourseSectionsQuery(
    { offeringId: offeringId || '' },
    { skip: !offeringId }
  );

  const sections: CourseSection[] = useMemo(
    () =>
      sectionsData?.sections
        .slice()
        .sort((a, b) => a.order - b.order) || [],
    [sectionsData]
  );

  const [fetchModules] = useLazyGetSectionModulesQuery();
  const [allModules, setAllModules] = useState<SectionModule[]>([]);
  const [modulesLoading, setModulesLoading] = useState(false);
  const [modulesError, setModulesError] = useState<any>(null);

  useEffect(() => {
    if (!sections.length) {
      setAllModules([]);
      return;
    }

    const fetchAllModules = async () => {
      setModulesLoading(true);
      setModulesError(null);
      try {
        const modulePromises = sections.map((section) =>
          fetchModules({ sectionId: section.id }).unwrap()
        );
        const results = await Promise.all(modulePromises);
        const combinedModules = results
          .flatMap((result) => result.modules)
          .sort((a, b) => a.order - b.order);
        setAllModules(combinedModules);
      } catch (err) {
        console.error('Failed to fetch modules:', err);
        setModulesError(err);
      } finally {
        setModulesLoading(false);
      }
    };

    fetchAllModules();
  }, [sections, fetchModules]);

  useEffect(() => {
    if (!sections.length) return;
    if (!selectedSectionId) {
      setSelectedSectionId(sections[0].id);
    }
  }, [sections, selectedSectionId]);

  const modules: SectionModule[] = useMemo(() => {
    if (!selectedSectionId) return [];
    return allModules
      .filter((m) => m.courseSectionId === selectedSectionId)
      .sort((a, b) => a.order - b.order);
  }, [allModules, selectedSectionId]);

  useEffect(() => {
    if (!modules.length) return;
    if (!selectedModuleId) {
      setSelectedModuleId(modules[0].id);
    } else if (selectedSectionId && selectedModuleId && !modules.some(m => m.id === selectedModuleId)) {
      setSelectedModuleId(modules[0]?.id || null);
    }
  }, [modules, selectedModuleId, selectedSectionId]);

  const activeModule = useMemo(() => {
    if (!selectedModuleId) return null;
    return allModules.find((m) => m.id === selectedModuleId) || null;
  }, [allModules, selectedModuleId]);
  
  const activeSection = sections.find((s) => s.id === selectedSectionId);

  const {
    data: meeting,
    isLoading: meetingLoading,
    error: meetingError,
  } = useGetZoomMeetingByModuleQuery(selectedModuleId || '', {
    skip: !selectedModuleId || !activeModule || activeModule.contentType !== 'zoom',
  });

  const {
    data: recordings,
    isLoading: recordingsLoading,
    error: recordingsError,
  } = useListZoomRecordingsQuery(meeting?.id || '', {
    skip: !meeting,
  });

  const handleBack = () => {
    navigate(ROUTES.STUDENT_DASHBOARD);
  };

  const handleSelectSection = (sectionId: string) => {
    setSelectedSectionId(sectionId);
    setSelectedModuleId(null);
    if (isMobile) {
      setSidebarOpen(false);
    }
  };

  const handleSelectModule = (moduleId: string) => {
    setSelectedModuleId(moduleId);
    if (isMobile) {
      setSidebarOpen(false);
    }
  };

  const handleJoinMeeting = () => {
    if (meeting?.joinUrl) {
      window.open(meeting.joinUrl, '_blank');
    }
  };

  const handleDownloadRecording = async (fileId: string) => {
    setDownloadingFileId(fileId);
    setDownloadError(null);
    try {
      await downloadFileWithAuth('zoom-recordings', fileId);
    } catch (err) {
      console.error('Failed to download recording:', err);
      const errorMsg = getErrorMessage(err, 'Failed to download recording');
      setDownloadError(errorMsg);
    } finally {
      setDownloadingFileId(null);
    }
  };

  const formatFileSize = (bytes?: number): string => {
    if (!bytes) return 'Unknown size';
    if (bytes < 1024) return `${bytes} B`;
    if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(2)} KB`;
    if (bytes < 1024 * 1024 * 1024) return `${(bytes / (1024 * 1024)).toFixed(2)} MB`;
    return `${(bytes / (1024 * 1024 * 1024)).toFixed(2)} GB`;
  };

  const handleLogout = useCallback(() => {
    dispatch(clearPublicAuth());
    dispatch(clearDashboardAuth());
    navigate(ROUTES.HOME);
  }, [dispatch, navigate]);

  if (authLoading || (!user && authLoading)) {
    return <Loading variant="fullscreen" />;
  }

  if (!user || user.role !== 'student') {
    return (
      <Box sx={{ p: 4, textAlign: 'center' }}>
        <Typography variant="h5" color="error">
          Access Denied. Student role required.
        </Typography>
      </Box>
    );
  }

  if (!offeringId) {
    return (
      <Box sx={{ p: 4, textAlign: 'center' }}>
        <Typography variant="h6">No class selected.</Typography>
      </Box>
    );
  }

  const anyLoading = offeringLoading || sectionsLoading || modulesLoading;

  return (
    <Box sx={{ display: 'flex', minHeight: '100vh', bgcolor: 'background.default' }}>
      <Box
        component="main"
        sx={{
          flexGrow: 1,
          display: 'flex',
          height: '100vh',
          overflow: 'hidden',
          width: '100%',
        }}
      >
        {isMobile ? (
          <Drawer
            variant="temporary"
            open={sidebarOpen}
            onClose={() => setSidebarOpen(false)}
            ModalProps={{ keepMounted: true }}
            sx={{
              '& .MuiDrawer-paper': { boxSizing: 'border-box', width: 280 },
            }}
          >
            <SidebarContent
              sections={sections}
              modules={allModules}
              selectedSectionId={selectedSectionId}
              selectedModuleId={selectedModuleId}
              onSelectSection={handleSelectSection}
              onSelectModule={handleSelectModule}
              sectionsLoading={sectionsLoading}
              modulesLoading={modulesLoading}
              sectionsError={sectionsError}
              modulesError={modulesError}
            />
          </Drawer>
        ) : (
          <Box
            component="nav"
            sx={{
              width: 280,
              borderRight: '1px solid',
              borderColor: 'divider',
              overflowY: 'auto',
              bgcolor: 'background.paper',
              flexShrink: 0,
            }}
          >
            <SidebarContent
              sections={sections}
              modules={allModules}
              selectedSectionId={selectedSectionId}
              selectedModuleId={selectedModuleId}
              onSelectSection={handleSelectSection}
              onSelectModule={handleSelectModule}
              sectionsLoading={sectionsLoading}
              modulesLoading={modulesLoading}
              sectionsError={sectionsError}
              modulesError={modulesError}
            />
          </Box>
        )}

        <Box
          sx={{
            flexGrow: 1,
            overflowY: 'auto',
            backgroundColor: 'background.default',
          }}
        >
          <Box sx={{ maxWidth: 1200, mx: 'auto', px: { xs: 2, sm: 3, md: 4 }, py: 4 }}>
            <Box
              sx={{
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'space-between',
                mb: 4,
                gap: 2,
              }}
            >
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                <IconButton
                  onClick={handleBack}
                  sx={{
                    border: '1px solid',
                    borderColor: 'divider',
                    '&:hover': {
                      bgcolor: 'action.hover',
                    },
                  }}
                >
                  <ArrowBackIcon />
                </IconButton>
                <Box>
                  <Typography variant="body2" color="text.secondary" sx={{ mb: 0.5 }}>
                    {offering?.courseName}
                  </Typography>
                  <Typography variant="h4" sx={{ fontWeight: 700, color: 'text.primary' }}>
                    {offering?.name || 'Class'}
                  </Typography>
                </Box>
              </Box>
              {isMobile && (
                <IconButton
                  onClick={() => setSidebarOpen(true)}
                  sx={{
                    border: '1px solid',
                    borderColor: 'divider',
                  }}
                >
                  <MenuIcon />
                </IconButton>
              )}
            </Box>

            {anyLoading && <Loading variant="skeleton" fullWidth />}

            {offeringError && (
              <Error
                message="Failed to load class information."
                onRetry={undefined}
                fullWidth
              />
            )}

            {!anyLoading && !offeringError && sections.length === 0 && (
              <Card sx={{ p: 4, textAlign: 'center' }}>
                <Typography variant="h6" color="text.secondary" sx={{ mb: 1 }}>
                  No curriculum available yet
                </Typography>
                <Typography variant="body2" color="text.secondary">
                  Course content will appear here once sections and modules are published.
                </Typography>
              </Card>
            )}

            {!anyLoading && !offeringError && sections.length > 0 && !activeModule && (
              <Card sx={{ p: 4, textAlign: 'center' }}>
                <Typography variant="h6" color="text.secondary" sx={{ mb: 1 }}>
                  Select a module to begin
                </Typography>
                <Typography variant="body2" color="text.secondary">
                  Choose a module from the sidebar to view its content.
                </Typography>
              </Card>
            )}

            {!anyLoading && !offeringError && activeModule && activeSection && (
              <Box sx={{ display: 'flex', flexDirection: 'column', gap: 3 }}>
                  <Card
                  sx={{
                    borderRadius: 3,
                    boxShadow: 2,
                    overflow: 'hidden',
                  }}
                >
                  <CardContent sx={{ p: 4 }}>
                    <Typography
                      variant="overline"
                      color="primary"
                      sx={{ fontSize: '0.75rem', fontWeight: 600, letterSpacing: 1 }}
                    >
                      Section
                    </Typography>
                    <Typography variant="h4" sx={{ fontWeight: 700, mt: 1, mb: 2 }}>
                      {activeSection.name}
                    </Typography>
                    {activeSection.description && (
                      <Typography variant="body1" color="text.secondary" sx={{ lineHeight: 1.7 }}>
                        {activeSection.description}
                      </Typography>
                    )}
                  </CardContent>
                </Card>

                <Card
                  sx={{
                    borderRadius: 3,
                    boxShadow: 2,
                    overflow: 'hidden',
                  }}
                >
                  <CardContent sx={{ p: 4 }}>
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mb: 2 }}>
                      <Typography
                        variant="overline"
                        color="primary"
                        sx={{ fontSize: '0.75rem', fontWeight: 600, letterSpacing: 1 }}
                      >
                        Module
                      </Typography>
                      <Chip
                        label={activeModule.contentType}
                        size="small"
                        color="primary"
                        variant="outlined"
                      />
                    </Box>
                    <Typography variant="h3" sx={{ fontWeight: 700, mb: 2 }}>
                      {activeModule.name}
                    </Typography>
                    {activeModule.description && (
                      <Typography variant="body1" color="text.secondary" sx={{ lineHeight: 1.7, mb: 3 }}>
                        {activeModule.description}
                      </Typography>
                    )}

                    {activeModule.contentType === 'zoom' && (
                      <Box sx={{ mt: 4 }}>
                        {meetingLoading && <Loading variant="skeleton" fullWidth />}
                        
                        {meetingError && (
                          <Alert severity="info" sx={{ mb: 3 }}>
                            Live session information is not available yet.
                          </Alert>
                        )}

                        {meeting && meeting.sectionModuleId === selectedModuleId && (
                          <Card
                            sx={{
                              bgcolor: 'primary.50',
                              border: '2px solid',
                              borderColor: 'primary.main',
                              borderRadius: 3,
                              overflow: 'hidden',
                              mb: 3,
                            }}
                          >
                            <CardContent sx={{ p: 4 }}>
                              <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mb: 3 }}>
                                <VideoCallIcon sx={{ fontSize: 40, color: 'primary.main' }} />
                                <Typography variant="h5" sx={{ fontWeight: 700 }}>
                                  Live Zoom Session
                                </Typography>
                              </Box>
                              
                              <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2, mb: 3 }}>
                                <Box>
                                  <Typography variant="body2" color="text.secondary" sx={{ mb: 0.5 }}>
                                    Topic
                                  </Typography>
                                  <Typography variant="body1" sx={{ fontWeight: 600 }}>
                                    {meeting.topic}
                                  </Typography>
                                </Box>
                                
                                {meeting.startTime && (
                                  <Box>
                                    <Typography variant="body2" color="text.secondary" sx={{ mb: 0.5 }}>
                                      Start Time
                                    </Typography>
                                    <Typography variant="body1" sx={{ fontWeight: 600 }}>
                                      {new Date(meeting.startTime).toLocaleString()}
                                    </Typography>
                                  </Box>
                                )}
                                
                                {meeting.duration && (
                                  <Box>
                                    <Typography variant="body2" color="text.secondary" sx={{ mb: 0.5 }}>
                                      Duration
                                    </Typography>
                                    <Typography variant="body1" sx={{ fontWeight: 600 }}>
                                      {meeting.duration} minutes
                                    </Typography>
                                  </Box>
                                )}
                                
                                {meeting.password && (
                                  <Box>
                                    <Typography variant="body2" color="text.secondary" sx={{ mb: 0.5 }}>
                                      Meeting Password
                                    </Typography>
                                    <Typography variant="body1" sx={{ fontWeight: 600, fontFamily: 'monospace' }}>
                                      {meeting.password}
                                    </Typography>
                                  </Box>
                                )}
                              </Box>

                              <MUIButton
                                variant="contained"
                                color="primary"
                                size="large"
                                startIcon={<VideoCallIcon />}
                                onClick={handleJoinMeeting}
                                sx={{
                                  py: 1.5,
                                  px: 4,
                                  fontSize: '1rem',
                                  fontWeight: 600,
                                  borderRadius: 2,
                                }}
                              >
                                Join Live Session
                              </MUIButton>
                            </CardContent>
                          </Card>
                        )}

                        {meeting && meeting.sectionModuleId === selectedModuleId && (
                          <>
                            {recordingsLoading && <Loading variant="skeleton" fullWidth />}
                            
                            {recordingsError && (
                              <Alert severity="warning" sx={{ mb: 3 }}>
                                Unable to load recordings. Please try again later.
                              </Alert>
                            )}

                            {recordings && recordings.length > 0 && (
                          <Box>
                            <Typography variant="h6" sx={{ fontWeight: 600, mb: 2 }}>
                              Session Recordings
                            </Typography>
                            <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
                              {recordings.map((rec) => (
                                <Card
                                  key={rec.id}
                                  sx={{
                                    borderRadius: 2,
                                    boxShadow: 1,
                                    '&:hover': {
                                      boxShadow: 3,
                                    },
                                    transition: 'all 0.2s',
                                  }}
                                >
                                  <CardContent>
                                    <Box
                                      sx={{
                                        display: 'flex',
                                        justifyContent: 'space-between',
                                        alignItems: 'center',
                                        gap: 2,
                                      }}
                                    >
                                      <Box sx={{ flexGrow: 1 }}>
                                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 1 }}>
                                          <PlayCircleOutlineIcon color="primary" />
                                          <Typography variant="subtitle1" sx={{ fontWeight: 600 }}>
                                            {rec.recordingType || 'Session Recording'}
                                          </Typography>
                                        </Box>
                                        {rec.recordingStartTime && (
                                          <Typography variant="body2" color="text.secondary" sx={{ mb: 0.5 }}>
                                            Recorded: {new Date(rec.recordingStartTime).toLocaleString()}
                                          </Typography>
                                        )}
                                        {rec.fileSize && (
                                          <Typography variant="body2" color="text.secondary">
                                            Size: {formatFileSize(rec.fileSize)}
                                          </Typography>
                                        )}
                                      </Box>
                                      <MUIButton
                                        variant="contained"
                                        color="primary"
                                        startIcon={
                                          downloadingFileId === rec.fileId ? (
                                            <CircularProgress size={16} color="inherit" />
                                          ) : (
                                            <DownloadIcon />
                                          )
                                        }
                                        onClick={() => handleDownloadRecording(rec.fileId)}
                                        disabled={downloadingFileId === rec.fileId}
                                        sx={{
                                          minWidth: 140,
                                        }}
                                      >
                                        {downloadingFileId === rec.fileId ? 'Downloading...' : 'Download'}
                                      </MUIButton>
                                    </Box>
                                  </CardContent>
                                </Card>
                              ))}
                            </Box>
                          </Box>
                            )}
                          </>
                        )}

                        {(!meeting || meeting.sectionModuleId !== selectedModuleId) && !meetingLoading && (
                          <Alert severity="info">
                            Live session details will appear here once the meeting is scheduled.
                          </Alert>
                        )}

                        {downloadError && (
                          <Alert severity="error" sx={{ mt: 2 }} onClose={() => setDownloadError(null)}>
                            {downloadError}
                          </Alert>
                        )}
                      </Box>
                    )}

                    {activeModule.contentType !== 'zoom' && (
                      <Alert severity="info">
                        Content for this module type ({activeModule.contentType}) will be available soon.
                      </Alert>
                    )}
                  </CardContent>
                </Card>
              </Box>
            )}
          </Box>
        </Box>
      </Box>
    </Box>
  );
}

interface SidebarContentProps {
  sections: CourseSection[];
  modules: SectionModule[];
  selectedSectionId: string | null;
  selectedModuleId: string | null;
  onSelectSection: (id: string) => void;
  onSelectModule: (id: string) => void;
  sectionsLoading: boolean;
  modulesLoading: boolean;
  sectionsError: any;
  modulesError: any;
}

function SidebarContent({
  sections,
  modules,
  selectedSectionId,
  selectedModuleId,
  onSelectSection,
  onSelectModule,
  sectionsLoading,
  modulesLoading,
  sectionsError,
  modulesError,
}: SidebarContentProps) {
  if (sectionsLoading || modulesLoading) {
    return (
      <Box sx={{ p: 3 }}>
        <Loading variant="skeleton" fullWidth />
      </Box>
    );
  }

  if (sectionsError || modulesError) {
    return (
      <Box sx={{ p: 3 }}>
        <Error
          message="Failed to load curriculum"
          onRetry={undefined}
          fullWidth
        />
      </Box>
    );
  }

  if (sections.length === 0) {
    return (
      <Box sx={{ p: 3 }}>
        <Typography variant="body2" color="text.secondary" sx={{ textAlign: 'center' }}>
          No sections available for this offering.
        </Typography>
      </Box>
    );
  }

  return (
    <Box sx={{ display: 'flex', flexDirection: 'column', height: '100%' }}>
      <Box sx={{ px: 3, py: 3, borderBottom: '1px solid', borderColor: 'divider' }}>
        <Typography variant="h6" sx={{ fontWeight: 700, color: 'text.primary' }}>
          Curriculum
        </Typography>
      </Box>
      <Box sx={{ flexGrow: 1, overflowY: 'auto' }}>
        <List disablePadding>
          {sections.map((section) => {
            const sectionModules = modules.filter(
              (m: SectionModule) => m.courseSectionId === section.id
            );
            const isSectionSelected = selectedSectionId === section.id;
            
            return (
              <Box key={section.id}>
                <ListItemButton
                  selected={isSectionSelected}
                  onClick={() => onSelectSection(section.id)}
                  sx={{
                    py: 1.5,
                    px: 3,
                    borderLeft: isSectionSelected ? '3px solid' : '3px solid transparent',
                    borderLeftColor: isSectionSelected ? 'primary.main' : 'transparent',
                    bgcolor: isSectionSelected ? 'action.selected' : 'transparent',
                    '&:hover': {
                      bgcolor: 'action.hover',
                    },
                  }}
                >
                  <ListItemText
                    primary={
                      <Typography
                        variant="subtitle2"
                        sx={{
                          fontWeight: isSectionSelected ? 700 : 600,
                          color: isSectionSelected ? 'primary.main' : 'text.primary',
                        }}
                      >
                        {section.name}
                      </Typography>
                    }
                  />
                </ListItemButton>
                {sectionModules.map((module: SectionModule) => {
                  const isModuleSelected = selectedModuleId === module.id;
                  return (
                    <ListItemButton
                      key={module.id}
                      selected={isModuleSelected}
                      onClick={() => {
                        onSelectSection(section.id);
                        onSelectModule(module.id);
                      }}
                      sx={{
                        pl: 5,
                        pr: 2,
                        py: 1,
                        borderLeft: isModuleSelected ? '3px solid' : '3px solid transparent',
                        borderLeftColor: isModuleSelected ? 'primary.main' : 'transparent',
                        bgcolor: isModuleSelected ? 'action.selected' : 'transparent',
                        '&:hover': {
                          bgcolor: 'action.hover',
                        },
                      }}
                    >
                      <ListItemText
                        primary={
                          <Typography
                            variant="body2"
                            sx={{
                              fontWeight: isModuleSelected ? 600 : 400,
                              color: isModuleSelected ? 'primary.main' : 'text.primary',
                            }}
                          >
                            {module.name}
                          </Typography>
                        }
                      />
                      <Chip
                        label={module.contentType}
                        size="small"
                        variant="outlined"
                        color={isModuleSelected ? 'primary' : 'default'}
                        sx={{ ml: 1 }}
                      />
                    </ListItemButton>
                  );
                })}
                {sectionModules.length === 0 && (
                  <Box sx={{ pl: 5, pr: 2, py: 1 }}>
                    <Typography variant="caption" color="text.secondary" sx={{ fontStyle: 'italic' }}>
                      No modules yet
                    </Typography>
                  </Box>
                )}
              </Box>
            );
          })}
        </List>
      </Box>
    </Box>
  );
}
