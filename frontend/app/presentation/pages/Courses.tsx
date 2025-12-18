import React, { useState } from 'react';
import { Box, Container } from '@mui/material';
import { useCourses } from '../hooks/useCourses';
import CourseCard from '../components/card/CourseCard';
import CourseSearchBar from '../features/courses/CourseSearchBar';
import CategorySideBar from '../features/courses/CategorySideBar';
import Pagination from '../features/courses/Pagination';
import Loading from '../components/common/Loading';
import Error from '../components/common/Error';
import { useNavigate } from 'react-router';
import { ROUTES } from '../../shared/constants/routes';

export default function Courses() {
  const navigate = useNavigate();
  const [selectedCategory, setSelectedCategory] = useState<string | undefined>();
  const [searchQuery, setSearchQuery] = useState<string>('');
  const [page, setPage] = useState(1);

  const pageSize = 12;

  const { courses, total, loading, error, refetch } = useCourses({
    searchQuery: searchQuery || undefined,
    category: selectedCategory,
    limit: pageSize,
    offset: (page - 1) * pageSize,
  });

  const handleSearch = (query: string) => {
    setSearchQuery(query);
    setPage(1);
  };

  const handleCategorySelect = (categoryId: string | undefined) => {
    setSelectedCategory(categoryId);
    setPage(1);
  };

  const handleViewDetail = (courseId: string) => {
    navigate(ROUTES.COURSE_DETAIL(courseId));
  };

  const handlePageChange = (newPage: number) => {
    setPage(newPage);
  };

  return (
    <Container maxWidth="xl" sx={{ py: 4, px: { xs: 2, md: 4 } }}>
      <Box sx={{ mb: 4, display: 'flex', justifyContent: 'center' }}>
        <Box sx={{ width: { xs: '100%', md: '70%' } }}>
          <CourseSearchBar onSearch={handleSearch} />
        </Box>
      </Box>

      <Box sx={{ display: 'flex', flexDirection: { xs: 'column', md: 'row' }, gap: 2 }}>
        <Box sx={{ width: { xs: '100%', md: '25%' }, p: 3 }}>
          <CategorySideBar
            onCategorySelect={handleCategorySelect}
            selectedCategory={selectedCategory}
          />
        </Box>

        <Box sx={{ width: { xs: '100%', md: '75%' } }}>
          {loading && <Loading variant="skeleton" fullWidth />}
          {error && (
            <Error
              message={error.message || 'Failed to load courses'}
              onRetry={() => refetch()}
              fullWidth
            />
          )}
          {!loading && !error && (
            <>
              <Box
                sx={{
                  display: 'grid',
                  gridTemplateColumns: {
                    xs: '1fr',
                    sm: 'repeat(2, 1fr)',
                    md: 'repeat(3, 1fr)',
                  },
                  gap: 3,
                }}
              >
                {courses.map((course) => (
                  <Box key={course.id}>
                    <CourseCard course={course} onViewDetail={handleViewDetail} />
                  </Box>
                ))}
              </Box>
              {courses.length === 0 && !loading && (
                <Box sx={{ textAlign: 'center', py: 8, color: 'text.secondary' }}>
                  No courses found
                </Box>
              )}
              <Box sx={{ display: 'flex', justifyContent: 'center', mt: 4 }}>
                <Pagination
                  currentPage={page}
                  pageSize={pageSize}
                  totalItems={total}
                  onPageChange={handlePageChange}
                />
              </Box>
            </>
          )}
        </Box>
      </Box>
    </Container>
  );
}

