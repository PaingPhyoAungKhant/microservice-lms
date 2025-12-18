import type { User } from '../../domain/entities/User';
import type { IUserRepository, UpdateUserData } from '../../domain/repositories/IUserRepository';

export interface UpdateUserInput {
  userId: string;
  username?: string;
  email?: string;
  role?: 'student' | 'instructor' | 'admin';
  status?: 'active' | 'inactive' | 'pending' | 'banned';
}

export interface UpdateUserOutput {
  user: User;
}

export class UpdateUserUseCase {
  constructor(private userRepository: IUserRepository) {}

  async execute(input: UpdateUserInput): Promise<UpdateUserOutput> {
    const data: UpdateUserData = {
      username: input.username,
      email: input.email,
      role: input.role,
      status: input.status,
    };

    Object.keys(data).forEach((key) => {
      if (data[key as keyof UpdateUserData] === undefined) {
        delete data[key as keyof UpdateUserData];
      }
    });

    const user = await this.userRepository.update(input.userId, data);

    return {
      user,
    };
  }
}

