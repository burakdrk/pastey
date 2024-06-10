//
//  HomeView.swift
//  pastey
//
//  Created by Burak Duruk on 2024-06-03.
//

import SwiftUI

struct HomeView: View {
    @EnvironmentObject var auth: AuthViewModel
    
    var body: some View {
        VStack {
            Text("Welcome, \(auth.user?.email ?? "User")")
                .font(.largeTitle)
                .fontWeight(.bold)
                .padding(.bottom, 30)
        }
    }
}

#Preview {
    HomeView()
        .environmentObject(AuthViewModel())
}
