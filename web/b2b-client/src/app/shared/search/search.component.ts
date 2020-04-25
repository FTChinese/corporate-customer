import { Component, Input } from '@angular/core';
import { Validators, FormControl, FormGroup } from '@angular/forms';
import { SearchService } from '../search.service';

@Component({
  selector: 'app-search',
  templateUrl: './search.component.html',
  styleUrls: ['./search.component.scss'],
})
export class SearchComponent {

  @Input() placeHolder = 'Enter a value to search';
  // @Output() submitted = new EventEmitter<AbstractControl>();

  form = new FormGroup({
    search: new FormControl('', Validators.required),
  });

  constructor(
    private service: SearchService,
  ) { }

  // Form submit handler. The FromControl is passed so that host
  // could perform further data process.
  onSubmit() {
    // this.submitted.emit(this.search);
    this.service.submitForm(this.form.get('search'));
  }
}
