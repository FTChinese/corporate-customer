import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FlashComponent } from './flash/flash.component';
import { FeedbackInvalidComponent } from './feedback-invalid/feedback-invalid.component';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MarkedPipe } from './marked.pipe';
import { SearchComponent } from './search/search.component';
import { DynamicFormComponent } from './dynamic-form/dynamic-form.component';
import { DynamicControlComponent } from './dynamic-control/dynamic-control.component';

@NgModule({
  declarations: [
    FlashComponent,
    FeedbackInvalidComponent,
    MarkedPipe,
    SearchComponent,
    DynamicFormComponent,
    DynamicControlComponent,
  ],
  imports: [
    CommonModule,
    ReactiveFormsModule,
    // FormsModule,
  ],
  exports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    FlashComponent,
    FeedbackInvalidComponent,
    MarkedPipe,
    SearchComponent,
    DynamicFormComponent,
  ]
})
export class SharedModule { }
