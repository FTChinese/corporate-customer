import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { DynamicControlComponent } from './dynamic-control.component';

describe('FormControlComponent', () => {
  let component: DynamicControlComponent;
  let fixture: ComponentFixture<DynamicControlComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ DynamicControlComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(DynamicControlComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
