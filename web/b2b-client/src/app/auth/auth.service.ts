import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { ILogin, StaffAccount, IProfile, IPasswords, IProfileForm, JWTAccount } from '../models/staff';
import { Observable, of } from 'rxjs';
import { tap, switchMap } from 'rxjs/operators';

@Injectable({
  providedIn: 'root'
})
export class AuthService {

  account: JWTAccount | null = null;
  redirectUrl: string;
  private storeKey = 'sy_user';

  get isLoggedIn(): boolean {
    if (!this.account) {
      const val = localStorage.getItem(this.storeKey);
      if (val) {
        this.account = JSON.parse(val);
      }
    }

    if (!this.account) {
      return false;
    }

    if (this.isExpired(this.account)) {
      return false;
    }
    return true;
  }

  private isExpired(account: JWTAccount): boolean {
    return (Date.now() / 1000) > account.expiresAt;
  }

  get authHeader(): HttpHeaders {
    return new HttpHeaders({
      Authorization: `Bearer ${this.account.token}`
    });
  }

  get displayName(): string {
    if (!this.account) {
      return '';
    }

    if (this.account.displayName) {
      return this.account.displayName;
    }

    return this.account.userName;
  }

  constructor(
    private http: HttpClient,
  ) {}

  login(credentials: ILogin): Observable<JWTAccount> {
    return this.http
      .post<JWTAccount>(
        '/api/login',
        credentials,
      )
      .pipe(
        tap(val => {
          this.account = val;
          localStorage.setItem(this.storeKey, JSON.stringify(val));
        })
      );
  }

  logout(): void {
    this.account = null;
    localStorage.removeItem(this.storeKey);
  }

  loadProfile(): Observable<IProfile> {
    return this.http.get<IProfile>(`/api/settings/profile`, {
      headers: this.authHeader
    });
  }

  updateProfile(formData: IProfileForm): Observable<boolean> {
    return this.http.patch<IProfileForm>(
        `/api/staff/${this.account.id}`,
        formData,
        {
          observe: 'response',
        }
      )
      .pipe(
        switchMap(resp => of(resp.status === 204))
      );
  }

  changePassword(pws: IPasswords): Observable<boolean> {
    return this.http.patch<IPasswords>(
        `/api/account/password`,
        pws,
        {
          observe: 'response',
          headers: this.authHeader,
        }
      )
      .pipe(
        switchMap(resp => of(resp.status === 204))
      );
  }
}
