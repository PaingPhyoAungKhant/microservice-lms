import type { IUserRepository } from '../../domain/repositories/IUserRepository';

export interface DeleteUserInput {
  userId: string;
}

export interface DeleteUserOutput {
  success: boolean;
}

export class DeleteUserUseCase {
  constructor(private userRepository: IUserRepository) {}

  async execute(input: DeleteUserInput): Promise<DeleteUserOutput> {
    await this.userRepository.delete(input.userId);

    return {
      success: true,
    };
  }
}

