import type { User } from '../../domain/entities/User';
import type { IUserRepository } from '../../domain/repositories/IUserRepository';

export interface GetUserInput {
  userId: string;
}

export interface GetUserOutput {
  user: User | null;
}

export class GetUserUseCase {
  constructor(private userRepository: IUserRepository) {}

  async execute(input: GetUserInput): Promise<GetUserOutput> {
    const user = await this.userRepository.findById(input.userId);

    return {
      user,
    };
  }
}

