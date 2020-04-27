import { Component} from '@angular/core';
import { FormBuilder, Validators } from '@angular/forms';
import { HttpErrorResponse } from '@angular/common/http';
import { Credentials, Passport } from '../../models/admin';
import { RequestError } from '../../models/request-result';
import { AuthService } from '../../core/auth.service';
import { Router, NavigationExtras } from '@angular/router';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss']
})
export class LoginComponent {

  loginForm = this.formBuilder.group({
    email: ['', [Validators.required, Validators.email]],
    password: ['', Validators.required],
  });


  formErr: Partial<Credentials> = {};
  errMsg: string;

  constructor(
    private formBuilder: FormBuilder,
    private authService: AuthService,
    private router: Router,
  ) { }

  onSubmit() {
    if (this.loginForm.invalid) {
      const nameErr = this.loginForm.getError('required', 'email');
      if (nameErr) {
        this.formErr.email = '请输入有效的邮箱地址';
      }

      const pwErr = this.loginForm.getError('required', 'password');
      if (pwErr) {
        this.formErr.password = '密码无效';
      }

      return;
    }

    this.loginForm.disable();

    this.authService
      .login(this.loginForm.value)
      .subscribe({
        next: (data: Passport) => {
          console.log(data);
          if (this.authService.isLoggedIn) {
            const redirect = this.authService.redirectUrl ? this.router.parseUrl(this.authService.redirectUrl) : '/';

            const navigationExtras: NavigationExtras = {
              queryParamsHandling: 'preserve',
              preserveFragment: true,
            };

            this.router.navigateByUrl(redirect, navigationExtras);
          }
        },
        error: (err: HttpErrorResponse) => {
          this.loginForm.enable();
          this.handleLoginError(err);
        },
      });
  }

  private handleLoginError(errResp: HttpErrorResponse) {
    console.log(errResp);

    const err = RequestError.fromResponse(errResp);

    console.log(err);

    if (err.notFound) {
      // this.loadFlash('Invalid credentials');
      this.errMsg = 'Invalid credentials';
      return;
    }

    if (err.invalid) {
      this.formErr = err.invalidObject;
      console.log(this.formErr);
      return;
    }

    // Fallback to any other errors.
    this.errMsg = err.message;
  }

  clearFeedback() {
    this.formErr = {};
  }
}
