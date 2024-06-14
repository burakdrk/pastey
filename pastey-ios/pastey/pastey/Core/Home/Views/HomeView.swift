//
//  HomeView.swift
//  pastey
//
//  Created by Burak Duruk on 2024-06-03.
//

import SwiftUI
import AlertKit

struct HomeView: View {
    @StateObject private var viewModel = HomeViewModel()
        
    var body: some View {
        NavigationStack {
            VStack {
                HeaderView
                
                ZStack {
                    List() {
                        ForEach(viewModel.entries) { entry in
                            EntryRowView(entry: entry)
                        }
                        .onDelete(perform: viewModel.deleteEntry)
                    }
                    .listStyle(.plain)
                    .environment(\.editMode, $viewModel.editMode)
                    .refreshable {
                        await viewModel.fetchEntries()
                    }
                    
                    if viewModel.entries.isEmpty {
                        Text("No entries found")
                            .foregroundColor(Color.theme.secondaryText)
                            .padding()
                            .frame(maxWidth: .infinity, alignment: .center)
                    }
                }
                
                BottomButtonsView
            }
            .navigationBarHidden(true)
        }
    }
}

#Preview {
    HomeView()
}

// MARK: - Header
extension HomeView {
    private var HeaderView: some View {
        HStack {
            Text("pastey")
                .font(.title)
                .bold()
            
            Spacer()
            
            Button(action: {
                withAnimation(.spring()) {
                    viewModel.editMode = viewModel.editMode.isEditing ? .inactive : .active
                }
            }, label: {
                Text(viewModel.editMode.isEditing ? "Done" : "Edit")
            })
        }
        .padding(.horizontal, 20)
        .padding(.vertical, 10)
    }
}

// MARK: - Bottom Buttons
extension HomeView {
    private var BottomButtonsView: some View {
        HStack {
            BigButton(text: "Copy", action: {
                viewModel.copy()
            }, isFetching: viewModel.isFetching)
            
            BigButton(text: "Paste", action: {
                viewModel.paste()
            }, isFetching: viewModel.isFetching)
        }
        .padding()
        .padding(.top, -2)
        .background(Color.theme.background)
    }
}
