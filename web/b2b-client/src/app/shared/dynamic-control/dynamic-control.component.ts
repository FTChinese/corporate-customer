import { Component, Input } from '@angular/core';
import { DynamicControl } from '../control';
import { FormGroup, AbstractControl } from '@angular/forms';

@Component({
  selector: 'app-dynamic-control',
  templateUrl: './dynamic-control.component.html',
  styleUrls: ['./dynamic-control.component.scss']
})
export class DynamicControlComponent {
  @Input() control: DynamicControl;
  @Input() form: FormGroup;

  // The current FormControl instance.
  get formControl(): AbstractControl {
    return this.form.get(this.control.key);
  }

  get isInvalid(): boolean {
    console.log('Checking validity');
    console.log(this.formControl.errors);
    return this.formControl.invalid && (this.formControl.dirty || this.formControl.touched);
  }

  // {
  //   min?: {
  //     min: number,
  //     actual: number
  //   };
  //   max?: {
  //     max: number,
  //     actual: number
  //   };
  //   required?: true;
  //   email?: true;
  //   minLength?: {
  //     requiredLength: number,
  //     acutualLength: number
  //   };
  //   maxLength?: {
  //     requiredLength: number,
  //     acutualLength: number
  //   };
  //   pattern?: {
  //     requiredPattern: string,
  //     actualValue: string
  //   };
  // API response errors:
  //   minssing: string,
  //   missing_field: string,
  //   invalid: string,
  //   already_exists: string,
  // }
  get errMsg(): string {
    const errors = this.formControl.errors;

    console.log(errors);

    if (errors.required) {
      return `${this.control.label} is required`;
    }

    if (errors.email) {
      return 'Please provide a valid email';
    }

    if (errors.min) {
      return 'Too low';
    }

    if (errors.max) {
      return 'Too high';
    }

    if (errors.minLength) {
      return 'Too short';
    }

    if (errors.maxLength) {
      return 'Too long';
    }

    if (errors.pattern) {
      return 'Pattern mismatched';
    }

    if (errors.missing) {
      return 'The requesting resource does not exist, or is removed';
    }

    if (errors.missing_field) {
      return `The value for ${this.control.label} is missing`;
    }

    if (errors.invalid) {
      return errors.invalid;
    }

    if (errors.already_exists) {
      return `The same value for ${this.control.label} already exists. Please use another one.`;
    }

    return '';
  }
}
