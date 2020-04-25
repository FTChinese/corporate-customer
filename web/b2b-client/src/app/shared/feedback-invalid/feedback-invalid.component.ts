import { Component, Input } from '@angular/core';

@Component({
  selector: 'app-feedback-invalid',
  template: `<div class="o-forms-input__error">{{message}}</div>`,
})
export class FeedbackInvalidComponent {
  @Input() message: string;
}
