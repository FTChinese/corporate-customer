import { Component, OnInit, Input } from '@angular/core';
import { DynamicControl } from '../control';
import { FormGroup } from '@angular/forms';
import { ControlService } from '../control.service';
import { FormService } from '../form.service';
import { Button } from '../button';

@Component({
  selector: 'app-dynamic-form',
  templateUrl: './dynamic-form.component.html',
  styleUrls: ['./dynamic-form.component.scss'],
  providers: [ControlService],
})
export class DynamicFormComponent implements OnInit {
  // Default button
  @Input() button: Button = {
    block: false,
    text: 'Save',
  };

  @Input() controls: DynamicControl[] = [];
  form: FormGroup;

  constructor(
    private controlService: ControlService,
    private formService: FormService,
  ) { }

  ngOnInit(): void {
    this.form = this.controlService.toFormGroup(this.controls);
    this.formService.errorReceived$.subscribe(reqErr => {
      console.log(reqErr);
    });
  }

  onSubmit() {
    const data = JSON.stringify(this.form.getRawValue());
    console.log(data);
    this.formService.submit(data);
  }

  setError() {
    console.log('Setting error manually');
    this.controls.forEach(ctrl => {
      this.form.get(ctrl.key).setErrors({
        already_exists: 'The same value already exists. Please use another one'
      });
    });
  }
}
