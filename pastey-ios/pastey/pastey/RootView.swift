//
//  RootView.swift
//  pastey
//
//  Created by Burak Duruk on 2024-06-06.
//

import SwiftUI

struct RootView: View {
    @EnvironmentObject var auth: AuthViewModel
    
    var body: some View {
        if auth.isLoggedIn {
            HomeView()
        } else {
            LoginView()
        }
    }
}

#Preview {
    RootView()
        .environmentObject(AuthViewModel())
}
