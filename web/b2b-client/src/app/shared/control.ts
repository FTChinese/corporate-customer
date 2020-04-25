import { ValidatorFn } from '@angular/forms';

type ControlType = 'textbox' | 'dropdown' | 'textarea';
type InputType = 'text' | 'email' | 'password' | 'number' | 'url' | 'search';

/**
 * @description Represent the <opiton> element.
 */
interface OptionElement {
  key: string;
  value: string;
}

interface ControlOptions {
  value?: any;
  key: string;
  validators?: ValidatorFn[];
  label?: string;
  desc?: string;
  placeholder?: string;
}

interface InputControlOptions extends ControlOptions {
  type: InputType;
}

interface DropdownControlOptions extends ControlOptions {
  options: OptionElement[];
}

/**
 * @description ControlBase configures a new FormGroup(controls: { [key: string]: AbstractControl })
 */
export class DynamicControl {
  // Config Angular FormControl
  value: any;    // new FromControl(value)
  key: string; // new FormGroup( {key: FormControl} ). Also used in label's for and input's formControlName attributes.
  validators?: ValidatorFn[]; // new FormControl(value, validators)

  // Override in suclass
  controlType: ControlType; // Use in the ngSwitch directive.

  // Config the attributes of HTML element.
  label: string;
  desc: string;
  placeholder: string;

  // Overridable
  type: InputType; // <input>'s type attribute.
  options: OptionElement[]; // <option> tag inside <select>.

  constructor(opts: ControlOptions) {
    this.value = opts.value || null;
    this.key = opts.key;
    this.validators = opts.validators;

    this.label = opts.label || '';
    this.desc = opts.desc || '';
    this.placeholder = opts.placeholder || '';
  }
}

export class InputControl extends DynamicControl {
  controlType: ControlType = 'textbox';
  type: InputType;

  constructor(opts: InputControlOptions) {
    super(opts);
    this.type = opts.type;
  }
}

export class DropdownControl extends DynamicControl {
  controlType: ControlType = 'dropdown';
  options: OptionElement[] = [];

  constructor(opts: DropdownControlOptions) {
    super(opts);
    this.options = opts.options;
  }
}

export class TextareaControl extends DynamicControl {
  controlType: ControlType = 'textarea';

  constructor(opts: ControlOptions) {
    super(opts);
  }
}
