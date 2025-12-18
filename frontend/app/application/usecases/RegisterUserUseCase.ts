import type { IAuthRepository, AuthResult } from '../../domain/repositories/IAuthRepository';
import { storage } from '../../infrastructure/storage/storage';

export interface RegisterUserInput {
  name: string;
  email: string;
  password: string;
  role?: 'student' | 'instructor';
}

export interface RegisterUserOutput {
  user: AuthResult['user'];
  tokens: AuthResult['tokens'];
}

export class RegisterUserUseCase {
  constructor(private authRepository: IAuthRepository) {}

  async execute(input: RegisterUserInput): Promise<RegisterUserOutput> {
    const result = await this.authRepository.register({
      name: input.name,
      email: input.email,
      password: input.password,
      role: input.role || 'student',
    });

    storage.setAccessToken(result.tokens.accessToken);
    if (result.tokens.refreshToken) {
      storage.setRefreshToken(result.tokens.refreshToken);
    }
    storage.setUser(result.user);

    return {
      user: result.user,
      tokens: result.tokens,
    };
  }
}

