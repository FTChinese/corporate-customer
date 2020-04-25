import { Component, OnInit } from '@angular/core';
import { navigation } from '../nav-item';

@Component({
  selector: 'app-sidebar',
  templateUrl: './sidebar.component.html',
  styleUrls: ['./sidebar.component.scss']
})
export class SidebarComponent implements OnInit {

  navItems = navigation;

  constructor() { }

  ngOnInit(): void {
  }

}
