class StorageService {
  private prefix = 'asto_lms_';

  private isBrowser(): boolean {
    return typeof window !== 'undefined';
  }

  private getKey(key: string): string {
    return `${this.prefix}${key}`;
  }

  /**
   * Validates that a token is a non-empty string
   * @param token - The token to validate
   * @param tokenType - The type of token (for logging purposes)
   * @returns true if token is valid, false otherwise
   */
  private validateToken(token: string | null | undefined, tokenType: string): boolean {
    if (!token || typeof token !== 'string') {
      if (import.meta.env.DEV) {
        console.warn(`[Storage] Invalid ${tokenType}: token is not a string`, { token, type: typeof token });
      }
      return false;
    }
    
    const trimmedToken = token.trim();
    if (trimmedToken.length === 0) {
      if (import.meta.env.DEV) {
        console.warn(`[Storage] Invalid ${tokenType}: token is empty or whitespace only`);
      }
      return false;
    }
    
    return true;
  }

  setItem(key: string, value: string): void {
    if (!this.isBrowser()) {
      return;
    }
    try {
      localStorage.setItem(this.getKey(key), value);
    } catch (error) {
      console.error('Error setting item in localStorage:', error);
    }
  }

  getItem(key: string): string | null {
    if (!this.isBrowser()) {
      return null;
    }
    try {
      return localStorage.getItem(this.getKey(key));
    } catch (error) {
      console.error('Error getting item from localStorage:', error);
      return null;
    }
  }

  removeItem(key: string): void {
    if (!this.isBrowser()) {
      return;
    }
    try {
      localStorage.removeItem(this.getKey(key));
    } catch (error) {
      console.error('Error removing item from localStorage:', error);
    }
  }

  clear(): void {
    if (!this.isBrowser()) {
      return;
    }
    try {
      const keys = Object.keys(localStorage);
      keys.forEach((key) => {
        if (key.startsWith(this.prefix)) {
          localStorage.removeItem(key);
        }
      });
    } catch (error) {
      console.error('Error clearing localStorage:', error);
    }
  }

  setAccessToken(token: string): void {
    if (!this.validateToken(token, 'access token')) {
      if (import.meta.env.DEV) {
        console.warn('[Storage] setAccessToken: Token validation failed, not storing');
      }
      return;
    }
    const trimmedToken = token.trim();
    this.setItem('accessToken', trimmedToken);
    if (import.meta.env.DEV) {
      console.debug('[Storage] setAccessToken: Token stored successfully', {
        tokenLength: trimmedToken.length,
        tokenPreview: trimmedToken.substring(0, 20) + '...',
      });
    }
  }

  getAccessToken(): string | null {
    const token = this.getItem('accessToken');
    if (import.meta.env.DEV) {
      console.debug('[Storage] getAccessToken:', {
        found: !!token,
        length: token?.length || 0,
        preview: token ? token.substring(0, 20) + '...' : null,
      });
    }
    if (token && !this.validateToken(token, 'access token')) {
      if (import.meta.env.DEV) {
        console.warn('[Storage] getAccessToken: Invalid token found, removing');
      }
      this.removeAccessToken();
      return null;
    }
    return token;
  }

  removeAccessToken(): void {
    this.removeItem('accessToken');
  }

  setRefreshToken(token: string): void {
    if (!this.validateToken(token, 'refresh token')) {
      return;
    }
    this.setItem('refreshToken', token.trim());
  }

  getRefreshToken(): string | null {
    const token = this.getItem('refreshToken');
    if (token && !this.validateToken(token, 'refresh token')) {
      this.removeRefreshToken();
      return null;
    }
    return token;
  }

  removeRefreshToken(): void {
    this.removeItem('refreshToken');
  }

  setUser(user: any): void {
    this.setItem('user', JSON.stringify(user));
  }

  getUser(): any | null {
    const userStr = this.getItem('user');
    if (userStr) {
      try {
        return JSON.parse(userStr);
      } catch {
        return null;
      }
    }
    return null;
  }

  removeUser(): void {
    this.removeItem('user');
  }

  setDashboardAccessToken(token: string): void {
    if (!this.validateToken(token, 'dashboard access token')) {
      if (import.meta.env.DEV) {
        console.warn('[Storage] setDashboardAccessToken: Token validation failed, not storing');
      }
      return;
    }
    const trimmedToken = token.trim();
    this.setItem('dashboard_accessToken', trimmedToken);
    if (import.meta.env.DEV) {
      console.debug('[Storage] setDashboardAccessToken: Token stored successfully', {
        tokenLength: trimmedToken.length,
        tokenPreview: trimmedToken.substring(0, 20) + '...',
      });
    }
  }

  getDashboardAccessToken(): string | null {
    const token = this.getItem('dashboard_accessToken');
    if (import.meta.env.DEV) {
      console.debug('[Storage] getDashboardAccessToken:', {
        found: !!token,
        length: token?.length || 0,
        preview: token ? token.substring(0, 20) + '...' : null,
      });
    }
    if (token && !this.validateToken(token, 'dashboard access token')) {
      if (import.meta.env.DEV) {
        console.warn('[Storage] getDashboardAccessToken: Invalid token found, removing');
      }
      this.removeDashboardAccessToken();
      return null;
    }
    return token;
  }

  removeDashboardAccessToken(): void {
    this.removeItem('dashboard_accessToken');
  }

  setDashboardRefreshToken(token: string): void {
    if (!this.validateToken(token, 'dashboard refresh token')) {
      return;
    }
    this.setItem('dashboard_refreshToken', token.trim());
  }

  getDashboardRefreshToken(): string | null {
    const token = this.getItem('dashboard_refreshToken');
    if (token && !this.validateToken(token, 'dashboard refresh token')) {
      this.removeDashboardRefreshToken();
      return null;
    }
    return token;
  }

  removeDashboardRefreshToken(): void {
    this.removeItem('dashboard_refreshToken');
  }

  setDashboardUser(user: any): void {
    this.setItem('dashboard_user', JSON.stringify(user));
  }

  getDashboardUser(): any | null {
    const userStr = this.getItem('dashboard_user');
    if (userStr) {
      try {
        return JSON.parse(userStr);
      } catch {
        return null;
      }
    }
    return null;
  }

  removeDashboardUser(): void {
    this.removeItem('dashboard_user');
  }

  clearDashboard(): void {
    this.removeDashboardAccessToken();
    this.removeDashboardRefreshToken();
    this.removeDashboardUser();
  }
}

export const storage = new StorageService();

