//
//  RootTabView.swift
//  pastey
//
//  Created by Burak Duruk on 2024-06-10.
//

import SwiftUI

struct RootTabView: View {
    @State private var selectedIndex = 2
    
    var body: some View {
        TabView(selection: $selectedIndex) {
            
            ProfileView()
                .tabItem {
                    Image(systemName: "person")
                    Text("Profile")
                }
                .tag(1)
            
            HomeView()
                .tabItem {
                    Image(systemName: "doc.on.clipboard")
                    Text("Clipboard")
                }
                .tag(2)
            
            SettingsView(selectedTab: $selectedIndex)
                .tabItem {
                    Image(systemName: "gear")
                    Text("Settings")
                }
                .tag(3)
            
        }
        .onAppear() {
            UITabBar.appearance().backgroundColor = UIColor(Color.theme.background)
        }
    }
}

#Preview {
    RootTabView()
        .environmentObject(AuthViewModel())
}
