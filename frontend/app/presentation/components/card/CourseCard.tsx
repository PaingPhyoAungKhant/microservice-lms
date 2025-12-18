import React from 'react';
import { Card, CardContent, CardMedia, Typography, Box, Chip } from '@mui/material';
import type { Course } from '../../../domain/entities/Course';
import { getFileDownloadUrl } from '../../../infrastructure/api/utils';

type Props = {
  course: Course;
  onViewDetail?: (courseId: string) => void;
};

export default function CourseCard({ course, onViewDetail }: Props) {
  if (!course) return null;

  const thumbnailUrl =
    course.thumbnailUrl ||
    (course.thumbnailId ? getFileDownloadUrl('course-thumbnails', course.thumbnailId) : '/placeholder-course.jpg');

  return (
    <Card
      sx={{
        width: '100%',
        height: '100%',
        display: 'flex',
        flexDirection: 'column',
        cursor: 'pointer',
        transition: 'all 0.2s',
        '&:hover': {
          opacity: 0.9,
          boxShadow: 3,
        },
      }}
      onClick={() => onViewDetail?.(course.id)}
    >
      <CardMedia
        component="div"
        sx={{
          aspectRatio: '1 / 1',
          width: '100%',
          backgroundImage: `url('${thumbnailUrl}')`,
          backgroundSize: 'cover',
          backgroundPosition: 'center',
          backgroundRepeat: 'no-repeat',
        }}
      />
      <CardContent
        sx={{
          display: 'flex',
          flexDirection: 'column',
          gap: 2,
          px: 2,
          py: 3,
          flexGrow: 1,
        }}
      >
        <Typography
          variant="h4"
          component="h4"
          sx={{
            color: 'primary.main',
            fontSize: '1.5rem',
            fontWeight: 600,
          }}
        >
          {course.name}
        </Typography>
        <Typography
          variant="body2"
          sx={{
            color: 'text.secondary',
            fontSize: '0.875rem',
            minHeight: '3em',
            display: '-webkit-box',
            WebkitLineClamp: 3,
            WebkitBoxOrient: 'vertical',
            overflow: 'hidden',
          }}
        >
          {course.description}
        </Typography>
        {course.categories && course.categories.length > 0 && (
          <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5, mt: 1 }}>
            {course.categories.slice(0, 3).map((category) => (
              <Chip
                key={category.id}
                label={category.name}
                size="small"
                color="primary"
                variant="outlined"
              />
            ))}
            {course.categories.length > 3 && (
              <Typography variant="caption" color="text.secondary">
                +{course.categories.length - 3} more
              </Typography>
            )}
          </Box>
        )}
      </CardContent>
    </Card>
  );
}

