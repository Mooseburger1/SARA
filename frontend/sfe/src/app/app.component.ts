import { Component, OnInit } from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {UserService} from './user.service';
@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent implements OnInit {
  
  title = "Answer to life is 42!";
  user = {email: "test@google.com"}

  constructor(private service: UserService) {}

  ngOnInit(): void {
    this.service.get().subscribe((data: any) => {
      console.log(data);
      this.user = data;
    });
  }

}