import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { PageNotFoundComponent } from './page-not-found/page-not-found.component';
import { CenterLayoutComponent } from './center-layout/center-layout.component';
import { HomeComponent } from './home/home.component';
import { AuthGuard } from '../core/guards/auth.guard';
import { MainLayoutComponent } from './main-layout/main-layout.component';

/**
 * Lazy loading every module.
 * The layout module acts like template inheritance in tranditional
 * server-side rendering engine.
 * @see https://medium.com/angular-in-depth/angular-routing-reusing-common-layout-for-pages-from-different-modules-440a23f86b57
 */
const routes: Routes = [
  {
    path: '',
    component: MainLayoutComponent,
    canActivate: [AuthGuard],
    children: [
      { path: '', component: HomeComponent },
    ]
  },
  {
    path: '',
    component: CenterLayoutComponent,
    children: [
      {
        path: 'login',
        loadChildren: () => import('../auth/auth.module').then(m => m.AuthModule),
      }
    ]
  },
  {
    path: '',
    component: CenterLayoutComponent,
    children: [
      {
        path: 'forgot-password',
        loadChildren: () => import('../forgot-password/forgot-password.module').then(m => m.ForgotPasswordModule),
      }
    ]
  },
  {
    path: '**',
    component: PageNotFoundComponent,
  }
];

@NgModule({
  imports: [
    RouterModule.forChild(routes)
  ],
  exports: [RouterModule]
})
export class LayoutRoutingModule { }
