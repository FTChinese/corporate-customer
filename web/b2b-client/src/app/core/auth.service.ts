import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable, of } from 'rxjs';
import { tap, switchMap } from 'rxjs/operators';
import { Credentials, Passport, Profile, Passwords } from '../models/admin';
import { apiUrl } from './endpoint';

@Injectable({
  providedIn: 'root'
})
export class AuthService {

  passport: Passport | null = null;
  redirectUrl: string;
  private storeKey = 'b2b_pp';

  get isLoggedIn(): boolean {
    // First try to see if account data is cached.
    if (!this.passport) {
      const val = localStorage.getItem(this.storeKey);
      if (val) {
        this.passport = JSON.parse(val);
      }
    }

    // If it is still missing after retrieving from cache.
    if (!this.passport) {
      return false;
    }

    if (this.isExpired(this.passport)) {
      this.logout();
      return false;
    }

    return true;
  }

  private isExpired(passport: Passport): boolean {
    return (Date.now() / 1000) > passport.expiresAt;
  }

  get authHeader(): HttpHeaders {
    return new HttpHeaders({
      Authorization: `Bearer ${this.passport.token}`
    });
  }

  get displayName(): string {
    if (!this.passport) {
      return '';
    }

    if (this.passport.displayName) {
      return this.passport.displayName;
    }

    return this.passport.email.split('@')[0];
  }

  constructor(
    private http: HttpClient,
  ) {}

  login(credentials: Credentials): Observable<Passport> {
    return this.http
      .post<Passport>(
        apiUrl.login,
        credentials,
      )
      .pipe(
        tap(val => {
          this.passport = val;
          localStorage.setItem(this.storeKey, JSON.stringify(val));
        })
      );
  }

  logout(): void {
    this.passport = null;
    localStorage.removeItem(this.storeKey);
  }

  loadProfile(): Observable<Profile> {
    return this.http.get<Profile>(apiUrl.profile, {
      headers: this.authHeader
    });
  }

  changePassword(pws: Passwords): Observable<boolean> {
    return this.http.patch<Passwords>(
        apiUrl.changePassword,
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
