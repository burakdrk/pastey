//
//  LoginView.swift
//  pastey
//
//  Created by Burak Duruk on 2024-06-04.
//

import SwiftUI

struct LoginView: View {
    @EnvironmentObject var auth: AuthViewModel
    @State private var email = ""
    @State private var password = ""
    
    var body: some View {
        VStack {
            Text("pastey")
                .font(.largeTitle)
                .fontWeight(.bold)
                .padding(.bottom, 30)
            
            VStack(spacing: 20) {
                TextField("Email", text: $email)
                    .padding()
                    .background(Color(.secondarySystemBackground))
                    .cornerRadius(10)
                    .autocapitalization(.none)
                    .keyboardType(.emailAddress)
                    .textContentType(.emailAddress)
                
                SecureField("Password", text: $password)
                    .padding()
                    .background(Color(.secondarySystemBackground))
                    .cornerRadius(10)
                
                BigButton(text: "Log in", action: {
                    auth.login(email: email, password: password)
                }, isFetching: auth.isFetching)
                                
                if let msg = auth.errorMessage {
                    Text(msg)
                        .font(.footnote)
                        .foregroundColor(.red)
                        .padding()
                }
            }.frame(maxWidth: 300)
            
            SignupView
        }
    }
}

extension LoginView {
    private var SignupView: some View {
        VStack {
            Text("Don't have an account?")
            Button("Sign up") {
                auth.signUp(email: email, password: password)
            }
        }.frame(height: 100)
    }
}

#Preview {
    LoginView()
        .environmentObject(AuthViewModel())
}
