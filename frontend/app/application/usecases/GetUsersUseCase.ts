import type {
  IUserRepository,
  UserQuery,
  UserListResult,
} from '../../domain/repositories/IUserRepository';

export interface GetUsersInput {
  searchQuery?: string;
  role?: 'student' | 'instructor' | 'admin';
  status?: 'active' | 'inactive' | 'pending' | 'banned';
  limit?: number;
  offset?: number;
  sortColumn?: string;
  sortDirection?: 'asc' | 'desc';
}

export interface GetUsersOutput {
  users: UserListResult['users'];
  total: number;
}

export class GetUsersUseCase {
  constructor(private userRepository: IUserRepository) {}

  async execute(input: GetUsersInput): Promise<GetUsersOutput> {
    const query: UserQuery = {
      searchQuery: input.searchQuery,
      role: input.role,
      status: input.status,
      limit: input.limit,
      offset: input.offset,
      sortColumn: input.sortColumn,
      sortDirection: input.sortDirection,
    };

    const result = await this.userRepository.find(query);

    return {
      users: result.users,
      total: result.total,
    };
  }
}

