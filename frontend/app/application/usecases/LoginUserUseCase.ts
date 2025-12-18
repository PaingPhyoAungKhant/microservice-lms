import type { IAuthRepository, AuthResult } from '../../domain/repositories/IAuthRepository';
import { storage } from '../../infrastructure/storage/storage';

export interface LoginUserInput {
  email: string;
  password: string;
}

export interface LoginUserOutput {
  user: AuthResult['user'];
  tokens: AuthResult['tokens'];
}

export class LoginUserUseCase {
  constructor(private authRepository: IAuthRepository) {}

  async execute(input: LoginUserInput): Promise<LoginUserOutput> {
    const result = await this.authRepository.login({
      email: input.email,
      password: input.password,
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

