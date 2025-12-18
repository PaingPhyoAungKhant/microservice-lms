import type { User } from '../../domain/entities/User';
import type { IUserRepository, CreateUserData } from '../../domain/repositories/IUserRepository';

export interface CreateUserInput {
  email: string;
  username: string;
  password: string;
  role: 'student' | 'instructor' | 'admin';
}

export interface CreateUserOutput {
  user: User;
}

export class CreateUserUseCase {
  constructor(private userRepository: IUserRepository) {}

  async execute(input: CreateUserInput): Promise<CreateUserOutput> {
    const data: CreateUserData = {
      email: input.email,
      username: input.username,
      password: input.password,
      role: input.role,
    };

    const user = await this.userRepository.create(data);

    return {
      user,
    };
  }
}

