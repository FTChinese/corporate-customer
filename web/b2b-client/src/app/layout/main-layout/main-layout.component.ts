import { Component } from '@angular/core';
import { AuthService } from 'src/app/auth/auth.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-main-layout',
  templateUrl: './main-layout.component.html',
  styleUrls: ['./main-layout.component.scss'],
})
export class MainLayoutComponent {

  constructor(
    public authService: AuthService,
    private route: Router,
  ) {
  }


  logout() {
    this.authService.logout();
    this.route.navigateByUrl('/login');
  }
}
