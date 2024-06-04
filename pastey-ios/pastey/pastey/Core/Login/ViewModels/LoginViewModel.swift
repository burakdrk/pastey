//
//  LoginViewModel.swift
//  pastey
//
//  Created by Burak Duruk on 2024-06-04.
//

import Foundation

class LoginViewModel: ObservableObject {
    @Published var email = ""
    @Published var password = ""
    @Published var errorMessage: String?
    
    init() {}
    
    func login() {
        guard !email.isEmpty, !password.isEmpty else {
            errorMessage = "Please fill in all fields"
            return
        }
        
        errorMessage = nil
        
        
        print("Log in")
    }
    
    func validate() {
        
    }
    
    
}
