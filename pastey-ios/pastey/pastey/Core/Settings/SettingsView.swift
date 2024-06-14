//
//  SettingsView.swift
//  pastey
//
//  Created by Burak Duruk on 2024-06-10.
//

import SwiftUI
import AlertKit

struct SettingsView: View {
    @State private var pickedColor = Color.theme.accent
    @State private var showingProfile = false
    @Binding var selectedTab: Int
    
    var body: some View {
        NavigationStack {
            VStack {
                GradientHeaderView(title: "Settings")
                
                Form {
                    Section(header: Text("Account")) {
                        Label("Profile", systemImage: "person").onTapGesture {
                            selectedTab = 1
                        }
                    }
                    
                    Section(header: Text("App")) {
                        Label("About", systemImage: "info.circle")
                        .onTapGesture {
                            showingProfile.toggle()
                        }.sheet(isPresented: $showingProfile) {
                            AboutView()
                        }
                    }
                    
                    Section(header: Text("Appearance (Must restart app)")) {
                        ColorPicker("Accent color", selection: $pickedColor, supportsOpacity: false)
                            .onChange(of: pickedColor) { _ in
                                UserDefaults.standard.set(pickedColor.toHex(), forKey: "accentColor")
                            }
                        
                        Button("Reset to default") {
                            UserDefaults.standard.removeObject(forKey: "accentColor")
                            AlertKitAPI.present(title: "Accent color reset", icon: .done, style: .iOS16AppleMusic, haptic: .success)
                        }
                    }
                }
                
                Spacer()
            }
        }.navigationBarHidden(true)
    }
}

#Preview {
    SettingsView(selectedTab: .constant(1))
}
