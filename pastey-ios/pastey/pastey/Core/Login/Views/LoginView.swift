//
//  LoginView.swift
//  pastey
//
//  Created by Burak Duruk on 2024-06-04.
//

import SwiftUI

struct LoginView: View {
    @StateObject private var viewModel = LoginViewModel()
    
    var body: some View {
        VStack {
            Text("pastey")
                .font(.largeTitle)
                .fontWeight(.bold)
                .padding(.bottom, 30)
            
            VStack(spacing: 20) {
                TextField("Email", text: $viewModel.email)
                    .padding()
                    .background(Color(.secondarySystemBackground))
                    .cornerRadius(10)
                    .autocapitalization(.none)
                    .keyboardType(.emailAddress)
                    .textContentType(.emailAddress)
                
                SecureField("Password", text: $viewModel.password)
                    .padding()
                    .background(Color(.secondarySystemBackground))
                    .cornerRadius(10)
                
                Button() {
                    viewModel.login()
                } label : {
                    Text("Log in")
                        .font(.headline)
                        .foregroundColor(.white)
                        .padding()
                        .frame(maxWidth: .infinity)
                        .background(Color.theme.accent)
                        .cornerRadius(10)
                }
                
                if let msg = viewModel.errorMessage {
                    Text(msg)
                        .font(.footnote)
                        .foregroundColor(.red)
                        .padding()
                }
            }.frame(maxWidth: 300)
                        
            VStack {
                Text("Don't have an account?")
                Button("Sign up") {
                    print("Sign up")
                }
            }.frame(height: 100)
        }
    }
}

#Preview {
    LoginView()
}
