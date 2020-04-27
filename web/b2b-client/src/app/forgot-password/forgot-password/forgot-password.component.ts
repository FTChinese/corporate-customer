import { Component, OnInit } from '@angular/core';
import { FormService } from 'src/app/shared/form.service';
import { Validators } from '@angular/forms';
import { InputControl, DynamicControl } from 'src/app/shared/control';
import { Button } from 'src/app/shared/button';

@Component({
  selector: 'app-forgot-password',
  templateUrl: './forgot-password.component.html',
  providers: [FormService],
})
export class ForgotPasswordComponent implements OnInit {

  dynamicControls: DynamicControl[] = [
    new InputControl({
      value: '',
      key: 'email',
      validators: [Validators.required, Validators.email, Validators.maxLength(64)],
      label: 'Email',
      placeholder: 'example@ftchinese.com',
      type: 'email',
    }),
  ];

  button: Button = {
    block: true,
    text: 'Send me a reset letter'
  };

  constructor(
    private formService: FormService,
  ) { }

  ngOnInit(): void {
    this.formService.formSubmitted$.subscribe(data => {
      console.log(data);
    });
  }

}
